import type { ChangeEvent, FormEvent } from "react";

import { DEFAULT_FORM_STATE, STATUS_OPTIONS } from "../constants";
import type { GameFormData } from "../types";

interface GameFormProps {
  formData: GameFormData;
  isEditing: boolean;
  isSubmitting: boolean;
  imagePreviewUrl: string;
  onChange: (event: ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => void;
  onFileChange: (event: ChangeEvent<HTMLInputElement>) => void;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
  onCancel: (nextState?: GameFormData) => void;
}

function GameForm({
  formData,
  isEditing,
  isSubmitting,
  imagePreviewUrl,
  onChange,
  onFileChange,
  onSubmit,
  onCancel,
}: GameFormProps) {
  return (
    <section className="panel editor-panel">
      <div className="panel-header">
        <div>
          <p className="section-kicker">Editor</p>
          <h2>{isEditing ? "Редактирование игры" : "Новая игра"}</h2>
          <p className="panel-text">Заполните карточку, загрузите постер и сохраните запись.</p>
        </div>
      </div>

      <form className="game-form" onSubmit={onSubmit}>
        <div className="image-upload-grid">
          <div className="upload-preview">
            {imagePreviewUrl ? (
              <img src={imagePreviewUrl} alt={formData.title || "Превью игры"} />
            ) : (
              <div className="upload-placeholder">
                <span>Game Cover</span>
                <strong>Загрузите постер</strong>
              </div>
            )}
          </div>

          <label className="upload-field">
            Обложка игры
            <input aria-label="Обложка игры" type="file" accept=".jpg,.jpeg,.png,.webp" onChange={onFileChange} />
            <span className="panel-text">JPG, PNG или WEBP. Путь сохранится в `imagePath`.</span>
          </label>
        </div>

        <label>
          Название
          <input name="title" placeholder="Например, The Last of Us Part II" value={formData.title} onChange={onChange} required />
        </label>

        <div className="form-row">
          <label>
            Жанр
            <input name="genre" placeholder="Action, RPG, Strategy" value={formData.genre} onChange={onChange} required />
          </label>

          <label>
            Платформа
            <input name="platform" placeholder="PC, PS5, Xbox, Switch" value={formData.platform} onChange={onChange} required />
          </label>
        </div>

        <div className="form-row">
          <label>
            Год выпуска
            <input
              name="releaseYear"
              type="number"
              min="1970"
              max={new Date().getFullYear() + 1}
              value={formData.releaseYear}
              onChange={onChange}
              required
            />
          </label>

          <label>
            Оценка
            <input name="rating" type="number" min="1" max="10" value={formData.rating} onChange={onChange} required />
          </label>
        </div>

        <label>
          Статус
          <select aria-label="Статус игры" name="status" value={formData.status} onChange={onChange}>
            {STATUS_OPTIONS.map((status) => (
              <option key={status.value} value={status.value}>
                {status.label}
              </option>
            ))}
          </select>
        </label>

        <label>
          Описание
          <textarea
            name="description"
            placeholder="Кратко опишите игру для большой карточки и каталога."
            value={formData.description}
            onChange={onChange}
            rows={4}
            required
          />
        </label>

        <label>
          Путь изображения
          <input name="imagePath" placeholder="/blob/..." value={formData.imagePath} onChange={onChange} readOnly />
        </label>

        <div className="form-actions">
          <button type="submit" className="primary-cta" disabled={isSubmitting}>
            {isSubmitting ? "Сохранение..." : isEditing ? "Обновить игру" : "Добавить игру"}
          </button>

          {isEditing && (
            <button type="button" className="ghost-button" onClick={() => onCancel(DEFAULT_FORM_STATE)}>
              Отмена
            </button>
          )}
        </div>
      </form>
    </section>
  );
}

export default GameForm;
