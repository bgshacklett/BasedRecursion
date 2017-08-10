#!/usr/bin/python3
import sys
import json
import itertools
from functools import reduce

new_a = dict()
counter = 0
mapper = dict()


def pathfinder(o, p):
    for i in p:
        o = o[i]
    return o


def up_counter():
    global counter
    counter += 1


# Set top levels for the object
def __init__(input_1, input_2):
    global new_a
    new_a["Changed"] = []
    new_a["Removed"] = []
    new_a["Added"] = []
    # Immediate check if it's the same and successfully exit.
    if input_1 == input_2:
        print("Everything matches")
        exit(0)
    # Lets check for top level key differences
    if set(input_1) == set(input_2):
        for k1 in input_1:
            if input_2[k1] != input_1[k1]:
                processor({k1: input_1[k1]}, {k1: input_2[k1]})
    else:
        for (k1, k2) in itertools.zip_longest(input_1, input_2):
            # Handle top level differences by checking index and just add in as 'new'
            try:
                list(input_1).index(k2)
            except ValueError:
                if k2 is not None:
                    new_a["Added"].append({"Path": [], "Key": k2, "Value": input_2[k2]})
            try:
                list(input_2).index(k1)
            except ValueError:
                if k1 is not None:
                    new_a["Removed"].append({"Path": [], "Key": k1, "Value": input_1[k1]})
        for k1 in input_1:
            try:
                if input_2[k1] != input_1[k1]:
                    try:
                        processor({k1: input_1[k1]}, {k1: input_2[k1]})
                    except TypeError as e:
                        raise e
            except KeyError:
                pass


def processor(input_1, input_2, path=None):
    global counter
    global mapper
    global new_a
    if path is None:
        path = []
    if input_1 != input_2:
        if isinstance(input_1, str):
            input_1 = [input_1]
        if isinstance(input_2, str):
            input_2 = [input_2]
        if len(list(input_1)) > 1:
            for (k1, k2) in itertools.zip_longest(input_1, input_2):
                proc = True  # This variable we use to control recursing in the event of a different key
                try:
                    list(input_1).index(k2)
                except ValueError:
                    if k2 is not None:
                        new_a["Added"].append({"Path": path, "Key": k2, "Value": input_2[k2]})
                        proc = False
                try:
                    list(input_2).index(k1)
                except ValueError:
                    if k1 is not None:
                        new_a["Removed"].append({"Path": path, "Key": k1, "Value": input_1[k1]})
                        proc = False
                if proc:
                    try:
                        processor({k1: input_1[k1]}, {k1: input_2[k1]}, path=path)
                    except TypeError:
                        print(input_1)
                        print(input_2)
            return
        for k1 in input_1:
            try:
                up_counter()
                if not isinstance(input_1, dict):
                    v1 = input_1
                else:
                    v1 = input_1[k1]
                if not isinstance(input_2, dict):
                    v2 = input_2
                else:
                    v2 = input_2[k1]
                if v1 != v2:
                    npath = path[:]

                    if not isinstance(v1, dict):
                        new_a["Changed"].append({"oldValue": v1, "newValue": v2, "Key": k1, "Path": path})
                    else:
                        npath.append(k1)
                        mapper[counter] = npath
                        processor(v1, v2, mapper[counter])
                        return

                processor(v1, v2, mapper[counter])
            except KeyError:
                pass


def format_out(parse_object):
    # I stole this function
    def merge(a, b, path=None):
        if path is None: path = []
        for key in b:
            if key in a:
                if isinstance(a[key], dict) and isinstance(b[key], dict):
                    merge(a[key], b[key], path + [str(key)])
                elif a[key] == b[key]:
                    pass
                else:
                    pass
                    #raise Exception('Conflict at %s' % '.'.join(path + [str(key)]))
            else:
                a[key] = b[key]
        return a
    return_object, return_added, return_removed, return_changed = dict(), dict(), dict(), dict()
    for k in parse_object["Changed"]:
        return_changed = [parse_object["Changed"][k]["oldValue"], parse_object["Changed"][k]["newValue"]]
        for p in reversed(parse_object["Changed"][k]["Path"]):
            return_changed = {p: return_changed}
        return_object = merge(return_object, return_changed)

    for k in parse_object["Removed"]:
        return_removed = {parse_object["Removed"][k]["Key"]: ["<"]}
        for p in reversed(parse_object["Removed"][k]["Path"]):
            return_removed = {p: return_removed}
        return_object = merge(return_object, return_changed)

    for k in parse_object["Added"]:
        return_added = {parse_object["Added"][k]["Key"]: [">"]}
        for p in reversed(parse_object["Added"][k]["Path"]):
            return_added = {p: return_added}
            return_object = merge(return_object, return_changed)
    return return_object #reduce(merge, [return_added, return_removed, return_changed])


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Only accepts two arguments.")
        exit(1)
    json_one = json.loads(open(sys.argv[1]).read())
    json_two = json.loads(open(sys.argv[2]).read())
    __init__(json_one, json_two)
    #print(json.dumps(format_out(new_a), sort_keys=False, indent=2))
    print(json.dumps(new_a, sort_keys=False, indent=2))



