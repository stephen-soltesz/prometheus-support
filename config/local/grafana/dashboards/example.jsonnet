local grafana = import "grafonnet/grafana.libsonnet";

grafana.dashboard.new(
    "Frontend Stats2",
    refresh="",
    timezone="utc",
    time_from="now-12h",
    time_to="now",
)
+ grafana.dashboard.addRow(
    grafana.row.new(
        id=2,
        height="500px"
    ) + grafana.row.addPanel(
        grafana.graphPanel.new(
            "Frontend QPS",
            id=3,
            span=6,
            format="ops",
            min=0,
            linewidth=2,
            datasource="Prometheus",
            legend_alignAsTable=true,
            legend_rightSide=true,
        ) + grafana.graphPanel.addTarget(
            grafana.prometheus.target(
                'sum by(code) (rate(http_requests_total{container="prometheus"}[2m]))',
                legendFormat="{{code}}",
            )
        )
    ) + grafana.row.addPanel(
        grafana.graphPanel.new(
            "Handler Latency",
            id=4,
            span=6,
            format="µs",
            min=0,
            linewidth=2,
            datasource="Prometheus",
            legend_alignAsTable=true,
            legend_rightSide=true,
        ) + grafana.graphPanel.addTarget(
            grafana.prometheus.target(
                'sum by (handler) (rate(http_request_duration_microseconds{quantile="0.9"}[2m]))',
                legendFormat="{{handler}}",
            )
        )
    )
)
