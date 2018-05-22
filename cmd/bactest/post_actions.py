from typing import List
from module import Module

POST_ACTION_SUFFIX = "post_"


class PostActions:
    def __init__(self):
        self.modules = []

    def __str__(self):
        return self.modules
    
    def __len__(self):
        return len(self.modules)