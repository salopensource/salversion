# Generated by Django 2.0.2 on 2018-02-27 09:16

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Plugin',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.TextField(blank=True, null=True)),
            ],
        ),
        migrations.CreateModel(
            name='Report',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('ip_address', models.GenericIPAddressField()),
                ('current_version', models.TextField(blank=True, null=True)),
                ('database', models.TextField(blank=True, null=True)),
                ('install_type', models.TextField(blank=True, null=True)),
                ('number_of_machines', models.TextField(blank=True, null=True)),
                ('date', models.DateTimeField(auto_now_add=True)),
            ],
            options={
                'ordering': ['-date'],
            },
        ),
        migrations.AddField(
            model_name='plugin',
            name='report',
            field=models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='collector.Report'),
        ),
    ]
