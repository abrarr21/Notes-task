import mongoose from "mongoose";

const noteSchema = new mongoose.Schema(
  {
    title: {
      type: String,
      trim: true,
      required: true,
    },

    description: {
      type: String,
      trim: true,
    },
  },
  { timestamps: true },
);

const notesModel = mongoose.model("notes", noteSchema);

export default notesModel;
