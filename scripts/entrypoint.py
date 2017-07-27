#!/usr/local/bin/python

import sys
import json
import requests


def convertToNumber(job):
    convertedJob = {}
    for k, v in job.iteritems():
        if isinstance(v, dict):
            convertedJob[k] = convertToNumber(v)
        elif '.number' in k:
            convertedJob[k.split('.')[0]] = int(v)
        else:
            convertedJob[k] = v
    return convertedJob


if __name__ == "__main__":
    config = json.loads(sys.argv[2]).get("vargs")
    server = config.get("server")
    job_config = convertToNumber(config.get("job_config"))
    print job_config

    r = requests.get(('%s/v1/jobs/%s') % (server, job_config.get("id")))
    if r.status_code == 200:
        r = requests.put('%s/v0/scheduled-jobs/%s' % (server, job_config.get("id")), json=job_config)
    elif r.status_code == 404:
        r = requests.post('%s/v0/scheduled-jobs/' % (server), json=job_config)

    print r.text
