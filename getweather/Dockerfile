FROM python:3.7.2-alpine3.8
ADD getweather.py /
ADD requirements.txt /
RUN pip install -r requirements.txt
CMD ["python", "/getweather.py"]
