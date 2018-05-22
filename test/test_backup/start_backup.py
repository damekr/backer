


class StartBackup:
    parallel = True
    def __init__(self):
        print("Hello from Start Backup")

    def run(self):
        print("Hello from start backup run")

class StartSecondBackup:
    parallel = False
    def __init__(self):
        print("Hello from Start Second Backup")

    def run(self):
        print("Hello from start second backup run")


def setup():
    print("SETUP Start backup")


def teardown():
    print("TearDown Start Backup")