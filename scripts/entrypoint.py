#!/usr/local/bin/python

import sys
import json
import requests
import collections

def update(d, u):
    for k, v in u.iteritems():
        if isinstance(v, collections.Mapping):
            r = update(d.get(k, {}), v)
            d[k] = r
        else:
            d[k] = u[k]
    return d

if __name__ == '__main__':
    config = json.loads(sys.argv[2]).get('vargs')
    server = config.get('server')
    job_config = config.get('job_config')
    print job_config

    r = requests.get(('%s/v1/jobs/%s') % (server, job_config.get('id')))
    if r.status_code == 200:
        print 'Updating job'
        job_config = update(job_config, json.loads(r.text))
        r = requests.put('%s/v0/scheduled-jobs/%s' % (server, job_config.get('id')), json=job_config)
    else:
        print 'Creating new job'
        r = requests.post('%s/v0/scheduled-jobs' % (server), json=job_config)

    print 'Status code:' + str(r.status_code)
    print 'Body: ' + r.text
