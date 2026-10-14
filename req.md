architecture

                    GitHub
                      │
                 GitHub Pages
                      │
                React + Vite
                      │
                HTTPS / API
                      ↓
                  Render
                      │
                Golang + Gin
                      │
          ┌───────────┴───────────┐
          │                       │
     /api/profile            /api/projects


req
1. data berasal dari req.md
2. ada interface khusus dengan password untuk crud data
3. bisa tambahkan image untuk ss project dan update photo profile
4. ada chatbot yg akan jawab pertanyaan recruiter bedasarkan data porto.md
5. chatbot harus bisa filter pertnyaan recruiter agar tidak melenceng dari porto.md
6. chatbot akan gunakan llm untuk menjawab / atau mengarahkan jawaban sesuai data porto.md
7. berikan interface untuk ganti llm
8. llm yg digunakan gemini / gpt / grook /claude
9. masukkan data porto.md masukkan ke db 
10. buat dokumen cara deployment github dan render