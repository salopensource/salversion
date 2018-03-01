from django.shortcuts import render
from django.http import HttpResponse
from .models import *
from django.views.decorators.csrf import csrf_exempt
from django.utils import timezone
import plistlib
import logging
import sys
import os
import requests

# Create your views here.


LOG_FILENAME = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 'salversion.log')

logging.basicConfig(filename=LOG_FILENAME, level=logging.INFO)

logger = logging.getLogger(__name__)

def update_current_version():
    req = requests.get('https://api.github.com/repos/salopensource/sal/releases')
    if req.status_code != requests.codes.ok:
        return

    version = req.json()[0]['tag_name']
    obj, created = Setting.objects.update_or_create(
        name='version',
        defaults={'value': version},
    )

    return obj

def current_version():
    if Setting.objects.filter(name='version').exists():
        now = timezone.now()
        version_object = Setting.objects.get(name='version')
        if version_object.updated < now - timezone.timedelta(days=1):
            print('Getting a new version')
            obj = update_current_version()
        else:
            obj = version_object
    else:
        obj = update_current_version()

    return obj.value

def get_client_ip(request):
    x_forwarded_for = request.META.get('HTTP_X_FORWARDED_FOR')
    if x_forwarded_for:
        ip = x_forwarded_for.split(',')[0]
    else:
        ip = request.META.get('REMOTE_ADDR')
    return ip

@csrf_exempt
def index(request):
    if request.method == 'POST':
        if 'data' in request.POST:
            try:
                data = plistlib.readPlistFromString(request.POST['data'])
            except Exception:
                # malformed data, just return
                e = sys.exc_info()[0]
                logger.info(e)
                HttpResponse(current_version())
            if 'version' in data:
                version = data['version']
            else:
                version = '2.4.0.538'
            try:
                report = Report(ip_address=get_client_ip(request), current_version=version,
                database=data['database'], number_of_machines=data['machines'], install_type=data['install_type'])
                report.save()
                for plugin in data['plugins']:
                    item = Plugin(report=report, name=plugin)
                    item.save()
            except:
                e = sys.exc_info()[0]
                logger.info(e)
                HttpResponse(current_version())


    return HttpResponse(current_version())
