import app from "./app.js";
import connectDB from "./config/database.js";
import notesRouter from "./routes/notes.route.js";

connectDB();

app.get("/", (req, res) => {
  res.send("Server running perfectly");
});

app.use("/api/notes", notesRouter);

app.listen(6969, () => {
  console.log("Server running on port: 6969");
});
