
from grafanalib.core import *


dashboard = Dashboard(
  title='Frontend Stats',
  time=Time('now-12h', 'now'),
  refresh='',
  rows=[
    Row(
      height=Pixels(500),
      panels=[
      Graph(
        title='Frontend QPS',
        dataSource='Prometheus',
        targets=[
          Target(
            expr='sum by(code) (rate(http_requests_total{container="prometheus"}[2m]))',
            legendFormat='{{code}}',
            refId='A',
          ),
        ],
        legend=Legend(
            alignAsTable=True,
            rightSide=True,
        ),
        yAxes=[
          YAxis(format=OPS_FORMAT),
          YAxis(format=SHORT_FORMAT),
        ],
      ),
      Graph(
        title='Frontend latency',
        dataSource='Prometheus',
        targets=[
          Target(
            expr='sum by(handler) (rate(http_request_duration_microseconds{quantile="0.9", container="prometheus"}[2m]))',
            legendFormat='{{handler}}',
            refId='A',
          ),
        ],
        legend=Legend(
            alignAsTable=True,
            rightSide=True,
        ),
        yAxes=single_y_axis(format=SECONDS_FORMAT),
      ),
    ]),
  ],
).auto_panel_ids()

