# AI Interview

## Cara Menjalankan

1. Clone repository ini

```
git clone git@github.com:fastcampus-backend-golang/ai-interview.git
```

2. Masuk ke direktori

```
cd ai-interview
```

3. Jalankan docker

```
make docker
```

4. Export environment variable

```
export OPENAI_API_KEY="YOUR_API_KEY"
export DB_URI="mongodb://username:password@localhost:27017/db?authSource=admin"
```

5. Jalankan aplikasi

```
make run
```

6. Buka browser dan akses `http://localhost:8080`

## Konten
- ai: client untuk mengakses API OpenAI
- data: client untuk MongoDB
- static: aset statis untuk halaman frontend
- page: halaman frontend
- model: model data untuk backend
- handler: handler di server backend