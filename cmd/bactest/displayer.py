import sys
import time


class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

class Framework:
    def __init__(self, output = sys.stderr, module_name: str = None):
        self.output = output
        self.module_name = module_name
        self.date_prefix = time.asctime()
        if module_name != None:
            self.main_prefix = "[BACTEST]" + "[{}]".format(self.date_prefix) + " [{}]".format(self.module_name)
        else:
            self.main_prefix = "[BACTEST]" + "[{}]".format(self.date_prefix)
        self.info_prefix = self.main_prefix + "[INFO] "
        self.error_prefix = self.main_prefix + "[ERROR] "
        self.warning_prefix = self.main_prefix + "[WARNING] "
        self.debug_prefix = self.main_prefix + "[DEBUG] "

    def print(self, value):
        print(bcolors.OKGREEN + self.main_prefix + value + bcolors.ENDC, file=self.output)
    
    def info(self, value):
        print(bcolors.OKGREEN + self.info_prefix + value + bcolors.ENDC, file=self.output)
    
    def error(self, value):
        print(bcolors.FAIL + self.error_prefix + value + bcolors.ENDC, file=self.output)

    def warning(self, value):
        print(bcolors.WARNING + self.warning_prefix + value + bcolors.ENDC, file=self.output)
    
    def debug(self, value):
        print(bcolors.OKBLUE +  self.debug_prefix + value + bcolors.ENDC, file=self.output)