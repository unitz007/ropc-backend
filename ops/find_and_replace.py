#!/bin/python3

# Read in the file
with open('latest-image.txt', 'r') as file:
    filedata = file.read()
    print(filedata)

# # Replace the target string
# filedata = filedata.replace('abcd', 'ram')
#
# # Write the file out again
# with open('file.txt', 'w') as file:
#     file.write(filedata)