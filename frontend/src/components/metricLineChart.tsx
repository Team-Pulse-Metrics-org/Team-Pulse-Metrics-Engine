import React from "react";
import Card from "../components/card";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

interface ChartDataPoint {
  label: string;
  value: number;
}

interface MetricLineChartProps {
  title: string;
  subtitle: string;
  data: ChartDataPoint[];
  color: string;
  valueLabel: string; // e.g., "pts", "issues"
  yAxisMax?: number;
  yAxisMin?: number;
}

export const MetricLineChart: React.FC<MetricLineChartProps> = ({
  title,
  subtitle,
  data,
  color,
  valueLabel,
  yAxisMax,
  yAxisMin = 0, // Defaults to 0 baseline unless specified
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
          <LineChart
            data={data}
            margin={{ top: 10, right: 5, left: 0, bottom: 5 }}
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
              domain={[
                yAxisMin,
                yAxisMax !== undefined ? yAxisMax : "dataMax + 20",
              ]}
              width={35}
            />

            <Tooltip
              cursor={{
                stroke: "#1b1d24",
                strokeWidth: 1,
                strokeDasharray: "4 4",
              }}
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

            <Line
              type="monotone"
              dataKey="value"
              stroke={color}
              strokeWidth={2.5}
              dot={false}
              activeDot={{ r: 5, strokeWidth: 0, fill: color }}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </Card>
  );
};
