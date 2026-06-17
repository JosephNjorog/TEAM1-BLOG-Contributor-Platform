import { useEffect } from "react";
import { useEditor, EditorContent, type Editor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import Link from "@tiptap/extension-link";
import Image from "@tiptap/extension-image";
import Placeholder from "@tiptap/extension-placeholder";

interface RichTextEditorProps {
  content: string;
  editable: boolean;
  onChange: (html: string, wordCount: number) => void;
}

function wordCount(editor: Editor): number {
  const text = editor.getText().trim();
  return text === "" ? 0 : text.split(/\s+/).length;
}

export function RichTextEditor({ content, editable, onChange }: RichTextEditorProps) {
  const editor = useEditor({
    extensions: [
      StarterKit,
      Link.configure({ openOnClick: false, autolink: true }),
      Image,
      Placeholder.configure({ placeholder: "Start writing your article..." }),
    ],
    content,
    editable,
    onUpdate: ({ editor }) => onChange(editor.getHTML(), wordCount(editor)),
  });

  useEffect(() => {
    if (editor) editor.setEditable(editable);
  }, [editor, editable]);

  if (!editor) return null;

  return (
    <div className="article-editor rounded-xl2 border border-surface-border bg-surface-card">
      {editable && <Toolbar editor={editor} />}
      <div className="px-5 py-4 text-zinc-100">
        <EditorContent editor={editor} />
      </div>
    </div>
  );
}

function Toolbar({ editor }: { editor: Editor }) {
  const btn = (active: boolean) =>
    `rounded-md px-2.5 py-1.5 text-sm font-medium ${active ? "bg-surface-raised text-zinc-100" : "text-zinc-400 hover:bg-surface-raised hover:text-zinc-100"}`;

  return (
    <div className="flex flex-wrap items-center gap-1 border-b border-surface-border px-3 py-2">
      <button type="button" className={btn(editor.isActive("bold"))} onClick={() => editor.chain().focus().toggleBold().run()}>
        B
      </button>
      <button type="button" className={btn(editor.isActive("italic")) + " italic"} onClick={() => editor.chain().focus().toggleItalic().run()}>
        I
      </button>
      <button type="button" className={btn(editor.isActive("heading", { level: 2 }))} onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}>
        H2
      </button>
      <button type="button" className={btn(editor.isActive("heading", { level: 3 }))} onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()}>
        H3
      </button>
      <button type="button" className={btn(editor.isActive("bulletList"))} onClick={() => editor.chain().focus().toggleBulletList().run()}>
        • List
      </button>
      <button type="button" className={btn(editor.isActive("orderedList"))} onClick={() => editor.chain().focus().toggleOrderedList().run()}>
        1. List
      </button>
      <button
        type="button"
        className={btn(editor.isActive("link"))}
        onClick={() => {
          const url = window.prompt("Link URL");
          if (url) editor.chain().focus().setLink({ href: url }).run();
        }}
      >
        Link
      </button>
      <button
        type="button"
        className={btn(false)}
        onClick={() => {
          const url = window.prompt("Image URL (Cloudinary upload widget lands in a later phase)");
          if (url) editor.chain().focus().setImage({ src: url }).run();
        }}
      >
        Image
      </button>
    </div>
  );
}
