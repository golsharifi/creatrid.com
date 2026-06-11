"use client";

/**
 * Client-side vault encryption (WebCrypto, AES-256-GCM).
 *
 * The passphrase never leaves the browser and Creatrid never stores it —
 * losing the passphrase means the file is unrecoverable, by design.
 *
 * Wire format: magic "CRV1" (4) | salt (16) | iv (12) | ciphertext.
 */

const MAGIC = new TextEncoder().encode("CRV1");
const PBKDF2_ITERATIONS = 310_000;

async function deriveKey(passphrase: string, salt: Uint8Array): Promise<CryptoKey> {
  const material = await crypto.subtle.importKey(
    "raw",
    new TextEncoder().encode(passphrase),
    "PBKDF2",
    false,
    ["deriveKey"]
  );
  return crypto.subtle.deriveKey(
    { name: "PBKDF2", salt: salt as BufferSource, iterations: PBKDF2_ITERATIONS, hash: "SHA-256" },
    material,
    { name: "AES-GCM", length: 256 },
    false,
    ["encrypt", "decrypt"]
  );
}

export async function encryptFile(file: File, passphrase: string): Promise<Blob> {
  const salt = crypto.getRandomValues(new Uint8Array(16));
  const iv = crypto.getRandomValues(new Uint8Array(12));
  const key = await deriveKey(passphrase, salt);
  const plaintext = await file.arrayBuffer();
  const ciphertext = await crypto.subtle.encrypt(
    { name: "AES-GCM", iv: iv as BufferSource },
    key,
    plaintext
  );
  return new Blob([MAGIC, salt, iv, new Uint8Array(ciphertext)], {
    type: "application/octet-stream",
  });
}

export async function decryptFile(data: ArrayBuffer, passphrase: string): Promise<ArrayBuffer> {
  const bytes = new Uint8Array(data);
  const magic = new TextDecoder().decode(bytes.slice(0, 4));
  if (magic !== "CRV1") {
    throw new Error("Not a Creatrid-encrypted file");
  }
  const salt = bytes.slice(4, 20);
  const iv = bytes.slice(20, 32);
  const ciphertext = bytes.slice(32);
  const key = await deriveKey(passphrase, salt);
  try {
    return await crypto.subtle.decrypt(
      { name: "AES-GCM", iv: iv as BufferSource },
      key,
      ciphertext as BufferSource
    );
  } catch {
    throw new Error("Wrong passphrase or corrupted file");
  }
}
