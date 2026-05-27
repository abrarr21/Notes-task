import app from "./app.js";
import connectDB from "./config/database.js";

connectDB();

app.get("/", (req, res) => {
  res.send("Server running perfectly");
});

app.listen(6969, () => {
  console.log("Server running on port: 6969");
});
