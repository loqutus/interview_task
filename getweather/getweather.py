#!/usr/bin/env python3

import pyowm
import os
import sys
from pyowm import OWM

SOURCE = 'openweathermap'

if __name__ == '__main__':
    API_KEY = os.getenv('OPENWEATHER_API_KEY')
    if not API_KEY:
        sys.exit('OPENWEATHER_API_KEY not defined')
    CITY_NAME = os.getenv('CITY_NAME')
    if not CITY_NAME:
        sys.exit('CITY_NAME not defined')
    owm = OWM(API_KEY)
    obs = owm.weather_at_place(CITY_NAME)
    w = obs.get_weather()
    print(
        'source={source}, city="{city}", description="{description}", temp={temp}, humidity={humidity}'.format(
            source=SOURCE, city=CITY_NAME,
            description=w.get_detailed_status(), temp=w.get_temperature(unit='celsius')['temp'],
            humidity=w.get_humidity()))
