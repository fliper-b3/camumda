import datetime
import json
import random
import requests
import time
from os import environ

camunda_url = 'http://localhost:8080/engine-rest/'
n_process_started = 5
quiet_time_s = 2

headers = {"Content-Type": "application/json"}


def formatDate(date):
    return ''.join(date.astimezone().isoformat(sep='T', timespec='milliseconds').rsplit(':', 1))


def getCurrentTime():
    return datetime.datetime.now()


def getCountInstancesStartedAfter(date):
    urlAdd = 'history/process-instance/count'
    data = {"startedAfter": formatDate(date)}
    return requests.post(camunda_url + urlAdd, data=json.dumps(data), headers=headers).json()


def getProcessDefinitions():
    urlAdd = 'process-definition'
    return requests.get(camunda_url + urlAdd, headers=headers).json()


def startRandomProcess(processList):
    processId = processList[random.randrange(len(processList))]['id']
    urlAdd = 'process-definition/' + processId + '/start'
    requestBody = {
        "variables": {
            "amount": {
                "value": "10000",
                "type": "Integer"
            },
            "creditor": {
                "value": "10000",
                "type": "Integer"
            },
            "invoiceNumber": {
                "value": "10000",
                "type": "Integer"
            },
            "invoiceCategory": {
                "value": "10000",
                "type": "Integer"
            },
            "invoiceClassification": {
                "value": "10000",
                "type": "Integer"
            }
        },
        "businessKey": "myBusinessKey"
    }
    return requests.post(camunda_url + urlAdd, data=json.dumps(requestBody), headers=headers).json()


def printInfo():
    print("""
        Camunda bpm process starter

        Camunda Host: %s
        N processes started: %d
        Quiet time: %d
    """ % (camunda_url, n_process_started, quiet_time_s))


def main():
    # Get environment variables
    global camunda_url, n_process_started, quiet_time_s

    if environ.get('CAMUNDA_HOST') is not None:
        camunda_url = 'http://' + environ.get('CAMUNDA_HOST') + '/engine-rest/'
    if environ.get('N_PROCESS_STARTED') is not None:
        n_process_started = int(environ.get('N_PROCESS_STARTED'))
    if environ.get('QUIET_TIME_S') is not None:
        quiet_time_s = int(environ.get('QUIET_TIME_S'))
    printInfo()

    # Start creating processes
    processDefinitions = getProcessDefinitions()
    while True:
        for i in range(n_process_started):
            print(startRandomProcess(processDefinitions))
        print(getCountInstancesStartedAfter(getCurrentTime() - datetime.timedelta(seconds=10)))
        time.sleep(quiet_time_s)


if __name__ == "__main__":
    main()
