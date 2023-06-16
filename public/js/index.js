const button = document.getElementById('btnRating');

  button.addEventListener('mousedown', function() {
    button.style.backgroundColor = 'red';
    // Ganti warna sesuai preferensi Anda
  });

  button.addEventListener('mouseup', function() {
    button.style.backgroundColor = 'blue';
    // Ganti warna kembali saat tombol dilepas
  });