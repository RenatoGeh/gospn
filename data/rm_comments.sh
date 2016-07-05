#!/bin/bash

sed -i.bak "/# CREATOR: GIMP PNM Filter Version 1.1/d" {$1}/*
