import app from "./app.js";
import connectDB from "./config/database.js";

connectDB();

app.get("/", (req, res) => {
  res.send("Server is running perfectly");
});

app.listen(6969, () => {
  console.log(`Server is running on port: 6969`);
});
