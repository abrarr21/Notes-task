import { Router } from "express";
import type { IRouter } from "express";
import notesModel from "../models/notes.model.js";

const notesRouter: IRouter = Router();

// @route /api/notes
// @description Create Note with title and description given by user
// @access Public
notesRouter.post("/", async (req, res) => {
  const { title, description } = req.body;

  if (!title) {
    return res.status(400).json({ message: "Title is required" });
  }

  if (title.trim().length < 3) {
    return res.status(400).json({
      message: "Title must at least 3 characters long",
    });
  }

  if (!description) {
    return res.status(400).json({
      message: "descrtiption is required",
    });
  }

  if (description.trim().length < 10) {
    return res.json(400).json({
      message: "descrtiption must be at least 10 characters long",
    });
  }

  const newNote = await notesModel.create({ title, description });

  return res.status(201).json({
    message: "Note created successfully",
    note: newNote,
  });
});

// @route /api/notes
// @description Return list of all notes
// @access Public
notesRouter.get("/", async (req, res) => {
  const allNotes = await notesModel.find();

  return res.status(200).json({
    message: "Notes fetched successfully",
    notes: allNotes,
  });
});

// @route /api/notes/:id
// @description Update description of the given note id
// @access Public
notesRouter.patch("/:id", async (req, res) => {
  const { id } = req.params;
  const { description } = req.body;

  if (!description) {
    return res.status(400).json({
      message: "descrtiption is required",
    });
  }

  if (description.trim().length < 10) {
    return res.status(400).json({
      message: "descrtiption must be at least 10 characters long",
    });
  }

  const noteExist = await notesModel.findById(id);
  if (!noteExist) {
    return res.status(204).json({ message: "Note not found" });
  }

  noteExist.description = description;
  noteExist.save();

  return res.status(200).json({
    message: "Note updated successfully",
    updatedNote: noteExist,
  });
});

export default notesRouter;
