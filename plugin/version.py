def rpmvercmp(a, b):
    """
    Compare alpha and numeric segments of two versions.
    return 1: a is newer than b
           0: a and b are the same version
          -1: b is newer than a

    This is based on the rpmvercmp function used in libalpm
    https://projects.archlinux.org/pacman.git/tree/lib/libalpm/version.c
    """
    if a == b:
        return 0

    one, two, ptr1, ptr2 = (0,)*4
    is_num = False

    # loop through each version segment of a and b and compare them
    while len(a) > one and len(b) > two:
        while len(a) > one and not a[one].isalnum():
            one += 1
        while len(b) > two and not b[two].isalnum():
            two += 1

        # if we ran to the end of either, we are finished with the loop
        if not (len(a) > one and len(b) > two):
            break

        # if the seperator lengths were different, we are also finished
        if one-ptr1 != two-ptr2:
            if one-ptr1 < two-ptr2:
                return -1
            else:
                return 1

        ptr1 = one
        ptr2 = two

        # grab first completely alpha or completely numeric segment leave one
        # and two pointing to the start of the alpha or numeric segment and
        # walk ptr1 and ptr2 to end of segment
        if a[ptr1].isdigit():
            while len(a) > ptr1 and a[ptr1].isdigit():
                ptr1 += 1
            while len(b) > ptr2 and b[ptr2].isdigit():
                ptr2 += 1
            is_num = True
        else:
            while len(a) > ptr1 and a[ptr1].isalpha():
                ptr1 += 1
            while len(b) > ptr2 and b[ptr2].isalpha():
                ptr2 += 1
            is_num = False

        # take care of the case where the two version segments are different
        # types: one numeric, the other alpha (i.e. empty) numeric segments are
        # always newer than alpha segments
        if two == ptr2:
            if is_num:
                return 1
            else:
                return -1


        if is_num:
            # we know this part of the strings only contains digits so we can
            # ignore the error value since it should always be nil
            a1 = int(a[one:ptr1])
            b1 = int(b[two:ptr2])

            # whichever number has more digits wins
            if a1 > b1:
                return 1

            if a1 < b1:
                return -1
        else:
            compare = alpha_compare(a[one:ptr1], b[two:ptr2])
            if compare < 0:
                return -1
            if compare > 0:
                return 1

        # advance one and two to next segment
        one = ptr1
        two = ptr2

    # this catches the case where all numeric and alpha segments have
    # compared identically but the segment separating characters were
    # different
    if len(a) <= one and len(b) <= two:
        return 0

    # the final showdown. we never want a remaining alpha string to beat an
    # empty string. the logic is a bit weird, but:
    # - if one is empty and two is not an alpha, two is newer.
    # - if one is an alpha, two is newer.
    # - otherwise one is newer.
    if (len(a) <= one and not b[two].isalpha()) or len(a) > one and a[one].isalpha():
        return -1

    return 1

def alpha_compare(a, b):
    """
    Compare two aplha strings a and b
    """
    if a == b:
        return 0

    i = 0
    while len(a) > i and len(b) > i and a[i] == b[i]:
        i += 1

    if len(a) == i and len(b) > i:
        return -1

    if len(b) == i:
        return 1

    return ord(a[i]) - ord(b[i])



# tests
alpha_numeric = [
    "1.0.1",
    "1.0.a",
    "1.0",
    "1.0rc",
    "1.0pre",
    "1.0p",
    "1.0beta",
    "1.0b",
    "1.0a",
]

numeric = [
    "20141130",
    "012",
    "11",
    "3.0.0",
    "2.011",
    "2.03",
    "2.0",
    "1.2",
    "1.1.1",
    "1.1",
    "1.0.1",
    "1.0.0.0.0.0",
    "1.0",
    "1",
]

git = [
    "r1000.b481c3c",
    "r37.e481c3c",
    "r36.f481c3c",
]

def bigger(versions):
    for i, v in enumerate(versions):
        for v2 in versions[i:]:
            if v != v2 and rpmvercmp(v, v2) != 1:
                print("failed: %s, %s" % (v, v2))

def smaller(versions):
    for i, v in reversed(list(enumerate(versions))):
        for v2 in versions[i:]:
            if v != v2 and rpmvercmp(v, v2) == -1:
                print("failed: %s, %s" % (v, v2))

bigger(alpha_numeric)
smaller(alpha_numeric)
bigger(numeric)
smaller(numeric)
bigger(git)
smaller(git)
