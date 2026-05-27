import app from "./app.js";
import connectDB from "./config/database.js";
import notesModel from "./models/notes.model.js";

connectDB();

app.get("/", (req, res) => {
  res.send("Server is running perfectly");
});

// @route /api/notes
// @title Title of the note
// @description Create description of the Note
// @access Public
app.post("/api/notes", async (req, res) => {
  const { title, description } = req.body;

  if (!title) {
    res.status(400).json({
      message: "Title is required",
    });
  }

  if (title.trim().length < 3) {
    res.status(400).json({
      message: "Title must be at least 3 characters long",
    });
  }

  if (!description) {
    res.status(400).json({
      message: "description is required",
    });
  }

  if (description.trim().length < 10) {
    res.status(400).json({
      message: "description must be at least 10 characters long",
    });
  }

  const newNote = await notesModel.create({ title, description });

  res.status(200).json({
    message: "Note created successfully",
    note: newNote,
  });
});

app.listen(6969, () => {
  console.log(`Server is running on port: 6969`);
});
