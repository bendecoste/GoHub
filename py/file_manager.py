class Manager:
    def __init__(self):
        self.files = {}

    def addFile(self, fname, contents):
        self.files[fname] = {"content" : contents}
