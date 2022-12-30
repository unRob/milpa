#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>

exec "$MILPA_COMPA" __doctor "${MILPA_OPT_SUMMARY:+--summary}"
