import React from "react";
import Card from "../components/card";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

interface ChartDataPoint {
  label: string; // X-Axis (e.g., "Jun 2")
  value: number; // Y-Axis metric number
}

interface MetricBarChartProps {
  title: string;
  subtitle: string;
  data: ChartDataPoint[];
  color: string;
  valueLabel: string; // e.g., "commits", "tasks resolved", "points"
  yAxisMax?: number;
}

export const MetricBarChart: React.FC<MetricBarChartProps> = ({
  title,
  subtitle,
  data,
  color,
  valueLabel,
  yAxisMax,
}) => {
  return (
    <Card className="p-6 bg-[#0d0e12] border-slate-800/80 rounded-2xl flex flex-col justify-between min-w-0">
      <div className="mb-6">
        <h3 className="text-base font-semibold tracking-tight text-slate-100">
          {title}
        </h3>
        <p className="text-xs text-slate-500 mt-0.5">{subtitle}</p>
      </div>

      <div className="h-64 w-full block">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart
            data={data}
            margin={{ top: 10, right: 5, left: 0, bottom: 5 }}
            barSize={38}
          >
            <CartesianGrid
              strokeDasharray="4 4"
              stroke="#1b1d24"
              vertical={false}
            />

            <XAxis
              dataKey="label"
              axisLine={false}
              tickLine={false}
              tick={{ fill: "#64748b", fontSize: 11, fontWeight: 500 }}
              dy={10}
            />

            <YAxis
              axisLine={false}
              tickLine={false}
              tick={{ fill: "#64748b", fontSize: 11, fontWeight: 500 }}
              domain={[0, yAxisMax !== undefined ? yAxisMax : "dataMax + 20"]}
              width={35}
            />

            <Tooltip
              cursor={{ fill: "rgba(27, 29, 36, 0.3)" }}
              formatter={(value: any) => [`${value} ${valueLabel}`, ""]}
              contentStyle={{
                backgroundColor: "#090a0f",
                borderColor: "#1b1d24",
                borderRadius: "12px",
                color: "#f8fafc",
              }}
              labelStyle={{
                color: "#64748b",
                fontSize: "11px",
                marginBottom: "4px",
              }}
              itemStyle={{ color: color, fontWeight: 600, fontSize: "13px" }}
            />

            <Bar dataKey="value" fill={color} radius={[4, 4, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </Card>
  );
};
