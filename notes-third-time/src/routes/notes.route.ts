import { Router } from "express";
import type { IRouter } from "express";
import notesModel from "../models/notes.model.js";

const notesRouter: IRouter = Router();

// @route /api/notes
// @description Create a note with title and descritpion given by user
// @access Public
notesRouter.post("/", async (req, res) => {
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
notesRouter.get("/", async (req, res) => {
  const allNotes = await notesModel.find();

  return res.status(200).json({
    message: "Notes fetched successfully",
    notes: allNotes,
  });
});

notesRouter.patch("/:id", async (req, res) => {
  const { id } = req.params;
  const { description } = req.body;

  if (!description) {
    return res.status(400).json({
      message: "description is required",
    });
  }

  if (description.trim().length < 10) {
    return res.status(400).json({
      message: "description must be atleast 10 characters long",
    });
  }

  const note = await notesModel.findById(id);
  if (!note) {
    return res.status(204).json({
      message: "Note not found",
    });
  }

  note.description = description;
  note.save();

  return res.status(200).json({
    message: "Note updated successfully",
    note: note,
  });
});

notesRouter.delete("/:id", async (req, res) => {
  const { id } = req.params;

  const note = await notesModel.findByIdAndDelete(id);
  if (!note) {
    return res.status(204).json({
      message: "Note not found",
    });
  }

  return res.status(200).json({
    message: "Note deleted successfully",
    deletedNote: note,
  });
});

export default notesRouter;
