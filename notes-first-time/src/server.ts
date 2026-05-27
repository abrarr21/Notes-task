import app from "./app.js";
import connectDB from "./config/database.js";
import notesModel from "./models/notes.model.js";

connectDB();

app.get("/", (req, res) => {
  res.send("Server is running perfectly");
});

// @route /api/notes
// @description Create title and description of the Note
// @access Public
app.post("/api/notes", async (req, res) => {
  const { title, description } = req.body;

  if (!title) {
    return res.status(400).json({
      message: "Title is required",
    });
  }

  if (title.trim().length < 3) {
    return res.status(400).json({
      message: "Title must be at least 3 characters long",
    });
  }

  if (!description) {
    return res.status(400).json({
      message: "description is required",
    });
  }

  if (description.trim().length < 10) {
    return res.status(400).json({
      message: "description must be at least 10 characters long",
    });
  }

  const newNote = await notesModel.create({ title, description });

  return res.status(201).json({
    message: "Note created successfully",
    note: newNote,
  });
});

// @route /api/notes
// @description Get all the notes
// @access Public
app.get("/api/notes", async (req, res) => {
  const notes = await notesModel.find();

  return res.status(200).json({
    message: "Notes fetched successfully",
    notes: notes,
  });
});

app.listen(6969, () => {
  console.log(`Server is running on port: 6969`);
});
