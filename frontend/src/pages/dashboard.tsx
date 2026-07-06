import Card from "../components/card";

function Dashboard(){
    return(
    <div className="bg-slate-950 p-6">
      <h1 className="text-2xl font-medium text-white pt-8 pl-5">
        Dashboard
      </h1>
      {/* stats section 4 cards */}
      <div className="grid grid-cols-4 gap-6 mt-8 px-6 h-42">
          <Card/>
          <Card/>
          <Card/>
          <Card/>
      </div>
      {/* chart section  */}
      <div className="grid grid-cols-3 gap-6 mt-8 px-6">
        <Card className="col-span-2 h-70"></Card>
        <Card className="col-span-1 h-70"></Card>
        <Card className="col-span-2 h-70"></Card>
        <Card className="col-span-1 h-70"></Card>
      </div>
      {/* Activity section */}
      <div className="grid grid-cols-4 mt-8 px-6 mb-8">
        <Card className="col-span-4 h-115"/>
      </div>
    </div>
    );
}

export default Dashboard;