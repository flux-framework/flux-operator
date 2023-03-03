import time
from functools import partial, update_wrapper


class Decorator:
    def __init__(self, func):
        update_wrapper(self, func)
        self.func = func

    def __get__(self, obj, objtype):
        return partial(self.__call__, obj)


class timed(Decorator):
    """
    Time the length of the run, add to times
    """

    def __call__(self, cls, *args, **kwargs):
        # Name of the key is after command
        if "timed_name" in kwargs:
            key = kwargs["timed_name"]
        # Fallback to name of function
        else:
            key = self.func.__name__

        start = time.time()
        res = self.func(cls, *args, **kwargs)
        end = time.time()
        cls.times[key] = round(end - start, 3)
        return res
