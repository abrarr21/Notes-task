import type { Request, Response } from "express";
import notesModel from "../models/notes.model.js";
import ApiResponse from "../utils/apiResponse.js";

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

  try {
    const newNote = await notesModel.create({ title, description });

    return res.status(201).json(
      new ApiResponse("Note created successfully", {
        newNote,
      }),
    );
  } catch (error) {
    console.log("error creating note", error);

    return res.status(500).json({
      message: "Internal Server Error",
    });
  }
};

export const getAllNotes = async (req: Request, res: Response) => {
  try {
    const allNotes = await notesModel.find();

    return res
      .status(200)
      .json(new ApiResponse("Notes fetched successfully", allNotes));
  } catch (error) {
    console.log("error fetching notes", error);
    return res.status(500).json({
      message: "Internal server error",
    });
  }
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

  try {
    const note = await notesModel.findById(id);
    if (!note) {
      return res.status(204).json({
        message: "Note not found",
      });
    }

    note.description = description;
    note.save();

    return res
      .status(200)
      .json(new ApiResponse("Note updated successfully", note));
  } catch (error) {
    console.log("error updating note", error);
    return res.status(500).json({
      message: "internal server error",
    });
  }
};

export const deleteNote = async (req: Request, res: Response) => {
  const { id } = req.params;

  try {
    const note = await notesModel.findByIdAndDelete(id);
    if (!note) {
      return res.status(204).json({
        message: "Note not found",
      });
    }

    return res
      .status(200)
      .json(new ApiResponse("Note deleted successfully", note));
  } catch (error) {
    console.log("error deleting note", error);
    return res.status(500).json({
      message: "internal server errror",
    });
  }
};
