import type { Request, Response } from "express";
import notesModel from "../models/notes.model.js";

export const createNote = async (req: Request, res: Response) => {
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
};

export const getAllNotes = async (req: Request, res: Response) => {
  const allNotes = await notesModel.find();

  return res.status(200).json({
    message: "Notes fetched successfully",
    notes: allNotes,
  });
};

export const updateNote = async (req: Request, res: Response) => {
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
};

export const deleteNote = async (req: Request, res: Response) => {
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
};
