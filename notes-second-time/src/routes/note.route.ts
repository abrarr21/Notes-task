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

export default notesRouter;
