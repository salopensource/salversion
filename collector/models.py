from django.db import models


class Report(models.Model):
    ip_address = models.GenericIPAddressField()
    current_version = models.TextField(null=True, blank=True)
    database = models.TextField(null=True, blank=True)
    install_type = models.TextField(null=True, blank=True)
    number_of_machines = models.TextField(null=True, blank=True)
    date = models.DateTimeField(auto_now_add=True)

    def __unicode__(self):
        return '%s %s' % (self.ip_address, self.date)
    
    class Meta:
        ordering = ['-date']


class Plugin(models.Model):
    report = models.ForeignKey(Report, on_delete=models.CASCADE)
    name = models.TextField(null=True, blank=True)


class Setting(models.Model):
    name = models.TextField(unique=True)
    value = models.TextField()
    updated = models.DateTimeField(auto_now=True)
