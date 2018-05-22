from typing import List

from module import Module

PRE_ACTION_SUFFIX = "pre_"


class PreActions:
    def __init__(self):
        self.modules = []

    def __str__(self):
        return self.modules
    
    def __len__(self):
        return len(self.modules)