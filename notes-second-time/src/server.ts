import app from "./app.js";
import connectDB from "./config/database.js";
import notesRouter from "./routes/note.route.js";

connectDB();

app.get("/", (req, res) => {
  res.send("server running perfectly");
});

app.use("/api/notes", notesRouter);

app.listen(6969, () => {
  console.log("Server running at port: 6969");
});
