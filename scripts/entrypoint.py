#!/usr/local/bin/python

import sys
import json
import requests

if __name__ == "__main__":
    config = json.loads(sys.argv[2]).get("vargs")
    server = config.get("server")
    job_config = config.get("job_config")
	
    r = requests.get('%s/v1/jobs/%s') % (server, job_config.get("id")))
    if r.status_code == 200:
        r = requests.put('%s/v0/scheduled-jobs/%s' % (server, job_config.get("id")), json=job_config)
    else if r.status_code == 404:
	r = reuests.post('%s/v0/scheduled-jobs/' % (server), json=job_config)
    
    print r.text
