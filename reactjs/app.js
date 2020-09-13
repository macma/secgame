const { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend } = Recharts;

const SimpleLineChart = React.createClass({
  render() {
    return (
      <LineChart width={600} height={300} data={data}
        margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
        <XAxis dataKey="id" />
        <YAxis label={{ value: 'Click(s)', angle: -90, position: 'insideLeft' }} />
        <CartesianGrid strokeDasharray="3 3" />
        <Tooltip />
        <Legend />
        <Line type="monotone" dataKey="black" stroke="#000000" activeDot={{ r: 8 }} />
        <Line type="monotone" dataKey="blue" stroke="#007aff" />
        <Line type="monotone" dataKey="orange" stroke="#ff9559" />
      </LineChart>
    );
  }
})