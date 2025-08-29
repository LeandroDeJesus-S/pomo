# ⏳ Pomo

Um timer de Pomodoro simples e eficiente feito em **Bash**, com notificações gráficas via `zenity`.  
Funciona no **Linux (Ubuntu 22.04+)** e roda em **background** até você parar manualmente.

## 📦 Instalação

```bash
git clone https://github.com/LeandroDeJesus-S/pomo.git
mkdir -p ~/.local/bin
cp pomo/pomo ~/.local/bin/pomo
chmod +x ~/.local/bin/pomo
```

Certifique-se de que `~/.local/bin` está no seu `$PATH`:

```bash
echo 'export PATH=$HOME/.local/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
```

O script usa **zenity** para notificações persistentes. Instale se não tiver:

```bash
sudo apt install zenity
```

---

## 🚀 Uso

### Iniciar um Pomodoro

```bash
# type "pomo" without args to get help
pomo <study:break[:long_break]>
```

* `<study>` → tempo de estudo em minutos
* `<break>` → tempo de pausa curta em minutos
* `<long_break>` → (opcional) tempo de pausa longa a cada 4 sessões (default = 15 min)

Exemplos:

```bash
pomo 25:5
# 25 min estudo, 5 min pausa curta, 15 min pausa longa (default)

pomo 30:10:20
# 30 min estudo, 10 min pausa curta, 20 min pausa longa
```

O Pomodoro **fica rodando em background até você parar** com `pomo stop`.

---

### Ver estatísticas

```bash
pomo stats
```

Exemplo de saída:

```
------ Pomodoro Stats ------
✅ Sessions finished: 3
⏳ Total study time: 75 min
📚 Current: Studying
⏱ Remaining: 12:34
----------------------------
```
