import csv
import json
import subprocess

# This script reads operations.csv file, which contains rows from Operations DB table i the format:
# "id","type","state","instance_id","created_at","updated_at","subaccountid","globalid"
# It groups operations with one subaccount and prints provisioning and deprovisioning operations for every subaccount.


class Operation:
    def __init__(self, row):
        self.id = row[0]
        self.type = row[1]
        self.state = row[2]
        self.instance_id = row[3]
        self.created_at = row[4]
        self.subaccountId = row[6]
        self.globalaccountId = row[7]


class Subaccount:
    def __init__(self, id, gid):
        self.id = id.strip('"')
        self.globalaccountId = gid
        self.operations = []

    def appendOperation(self, op : Operation):
        self.operations.append(op)
        self.operations.sort(key=lambda x: x.created_at)

    def print(self):
        print("Subaccount '"+self.id+"' (GID="+self.globalaccountId+"):")
        lastType = ""
        lastInstanceId = ""
        existingInstances = {}
        lastDeprovision = ""
        for op in self.operations:
            if op.instance_id != lastInstanceId or op.type != lastType:
                print("    " + op.created_at + "\t" + op.instance_id + "\t" + op.type + "\t" + op.id)
                lastType = op.type
                lastInstanceId = op.instance_id
                if op.type == "provision":
                    existingInstances[op.instance_id] = op.created_at
                elif op.instance_id in existingInstances:
                    del existingInstances[op.instance_id]
                    # lastDeprovision = op.created_at
        print("    Still exists: " + str(existingInstances))
        # print("    LastDeprovision: " + lastDeprovision)


subaccounts = {}

with open("operations.csv") as csvfile:
    r = csv.reader(csvfile)
    for row in r:
        op = Operation(row)
        if op.subaccountId not in subaccounts:
            subaccounts[op.subaccountId] = Subaccount(op.subaccountId, op.globalaccountId)
        subaccounts[op.subaccountId].appendOperation(op)

for v in subaccounts:
    subaccounts[v].print()

# out = open("out.sh")
for id in subaccounts:
    sid = id.strip('"')
    # print("------")
    if sid == 'subid':
        continue
    # res = subprocess.run('/Users/i321040/go/src/github.com/kyma-project/control-plane/tools/cli/./kcp-darwin', ['rt', '--help'],
    #                 {"KCPCONFIG":"/Users/i321040/Downloads/Downloads/kcp-prod.yaml", "HOME":"/Users/i321040"} )
    # os.execve("echo", ['$KCPCONFIG'], {"KCPCONFIG":"/Users/i321040/Downloads/Downloads/kcp-prod.yaml", "HOME":"/Users/i321040"})

    result = subprocess.run(['/Users/i321040/go/src/github.com/kyma-project/control-plane/tools/cli/./kcp-darwin', "rt", "-s", sid, "-o", "json"],
                   stdout=subprocess.PIPE, text=True, env = {"KCPCONFIG":"/Users/i321040/Downloads/Downloads/kcp-prod.yaml", "HOME":"/Users/i321040"})
    obj = json.loads(result.stdout)
    plaftormRegion = obj['data'][0]['subAccountRegion']
    gid = obj['data'][0]['globalAccountID']
    plan = obj['data'][0]['servicePlanName']
    if plan == "free":
        planType = 'free'
    elif plan == 'azure_lite':
        planType = 'tdd'
    else:
        planType = 'standard'
    # print(v, plaftormRegion, gid, plan)
    # print("\n\n# SID: " + sid + " GID: " + gid + " plan: "
    #       + plan + " plantype: " + planType + " region: "+ plaftormRegion)
    # print("echo \"registering SID=" + sid + "\"")
    # print("./edp register " + sid + " " + plaftormRegion + " " + planType)
    # print("./edp get " + sid)
    print(gid + " " + sid + " " + plan + " " + plaftormRegion + " " + obj['data'][0]['userID'])
