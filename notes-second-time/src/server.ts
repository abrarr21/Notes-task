import app from "./app.js";

app.get("/", (req, res) => {
  res.send("server running perfectly");
});

app.listen(6969, () => {
  console.log("Server running at port: 6969");
});
