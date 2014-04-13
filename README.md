docclean
========

[![Build Status](https://secure.travis-ci.org/StefanKjartansson/docclean.png)](http://travis-ci.org/StefanKjartansson/docclean)

Simple utility to strip out empty comments from python files.

### Before

```python
def foo():
    """
    """
    pass

def foo():
    '''
    '''
    pass

def foo():
    """
    Undocumented.
    """
    pass

def foo():
    """
    Undocumented. moar
    """
    pass
```

### After

```python
def foo():
    pass

def foo():
    pass

def foo():
    pass

def foo():
    """
    Undocumented. moar
    """
    pass
```

