CREATE TABLE notification_settings (
    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    setting_key VARCHAR(50) UNIQUE NOT NULL,
    setting_value TEXT NOT NULL,
    description VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO notification_settings (setting_key, setting_value, description) VALUES
('wa_check_start', '08:00', 'Jam mulai pengecekan notifikasi WA (format HH:MM)'),
('wa_check_end', '09:00', 'Jam akhir pengecekan notifikasi WA (format HH:MM)'),
('wa_enabled', 'true', 'Aktifkan/nonaktifkan notifikasi WA (true/false)'),
('wa_message_template', 'Assalamualaikum, kami informasikan bahwa anak Bapak/Ibu *{nama}* (NISN: {nisn}, Kelas: {kelas}) hari ini tercatat *{status}*. Mohon perhatiannya. Terima kasih.', 'Template pesan WA');
