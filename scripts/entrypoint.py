#!/usr/local/bin/python

import sys
import json
import requests

if __name__ == '__main__':
    config = json.loads(sys.argv[2]).get('vargs')
    server = config.get('server')
    job_config = config.get('job_config')
    print job_config

    r = requests.get(('%s/v1/jobs/%s') % (server, job_config.get('id')))
    if r.status_code == 200:
        print 'Updating job'
        updated_job = json.loads(r.text)
        updated_job.update(job_config)
        r = requests.put('%s/v0/scheduled-jobs/%s' % (server, updated_job.get('id')), json=updated_job)
    else:
        print 'Creating new job'
        r = requests.post('%s/v0/scheduled-jobs' % (server), json=job_config)

    print 'Status code:' + str(r.status_code)
    print 'Body: ' + r.text
