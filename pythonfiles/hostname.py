from powerline_shell.utils import BasicSegment
from os import uname


class Segment(BasicSegment):
    def add_to_powerline(self):

        hostname = uname().nodename
        if hostname != "sandbox":
            return
        bg = 120
        fg = 18
        self.powerline.append(" " + hostname + " ", fg, bg)
