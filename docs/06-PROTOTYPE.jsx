import { useState, useEffect } from "react";

const C = {
  primary: "#0F8B4F", primaryLight: "#E8F5EE", primaryDark: "#0A6B3C",
  accent: "#F5A623", bg: "#FAFBFC", card: "#FFFFFF", text: "#1A1D21",
  sub: "#6B7280", muted: "#9CA3AF", border: "#E5E7EB", borderL: "#F3F4F6",
  green: "#22C55E", yellow: "#EAB308", blue: "#3B82F6", red: "#EF4444",
};

const BADGES = [
  { id: "first_win", icon: "🏅", name: "Первая победа", earned: true, p: 100 },
  { id: "ten_wins", icon: "🎖", name: "Десятка", earned: true, p: 100 },
  { id: "fifty", icon: "⭐", name: "Полтинник", earned: false, p: 68 },
  { id: "streak", icon: "🔥", name: "Серия", earned: false, p: 60 },
  { id: "tourney", icon: "🏆", name: "Турнирный", earned: false, p: 40 },
  { id: "lvlup", icon: "👑", name: "Уровень вверх", earned: true, p: 100 },
  { id: "social", icon: "🤝", name: "Социальный", earned: true, p: 100 },
  { id: "vet", icon: "🎾", name: "Ветеран", earned: false, p: 34 },
  { id: "pop", icon: "👥", name: "Авторитет", earned: false, p: 70 },
  { id: "reg", icon: "📅", name: "Регулярный", earned: false, p: 75 },
];

const EVENTS = [
  { id: 1, title: "Ищу партнёра на вечер", type: "find_partner", sc: C.green, level: "3.0-4.0", time: "Сегодня, 18:00", loc: "NTC Astana", spots: "1/2", comp: "Singles", creator: "Алексей М." },
  { id: 2, title: "Парная игра выходного дня", type: "organized_game", sc: C.yellow, level: "2.5-3.5", time: "Сб, 10:00", loc: "Tennis Club Astana", spots: "6/8", comp: "Doubles", creator: "NTC Astana" },
  { id: 3, title: "Весенний турнир NTC", type: "tournament", sc: C.blue, level: "3.0+", time: "15-16 Мар", loc: "NTC Astana", spots: "12/16", comp: "Singles", creator: "NTC Astana" },
  { id: 4, title: "Тренировка с тренером", type: "training", sc: C.green, level: "1.0-2.5", time: "Чт, 17:00", loc: "Mega Tennis", spots: "3/6", comp: "Group", creator: "Mega Tennis" },
];

const PLAYERS = [
  { id: 1, name: "Марат Касымов", lv: 4.5, r: 1650, dist: "Есильский", g: "М", games: 120, wr: 72, av: "МК", on: true },
  { id: 2, name: "Алия Нурланова", lv: 3.5, r: 1400, dist: "Сарыаркинский", g: "Ж", games: 85, wr: 61, av: "АН", on: false },
  { id: 3, name: "Дмитрий Ли", lv: 4.0, r: 1520, dist: "Есильский", g: "М", games: 95, wr: 65, av: "ДЛ", on: true },
  { id: 4, name: "Сауле Абдрахманова", lv: 3.0, r: 1280, dist: "Алматинский", g: "Ж", games: 42, wr: 55, av: "СА", on: false },
  { id: 5, name: "Арман Жумабаев", lv: 3.5, r: 1380, dist: "Есильский", g: "М", games: 68, wr: 59, av: "АЖ", on: true },
  { id: 6, name: "Тимур Сагинтаев", lv: 2.5, r: 1150, dist: "Сарыаркинский", g: "М", games: 28, wr: 50, av: "ТС", on: false },
];

const MSGS = [
  { id: 1, text: "Привет! Играем завтра в 18?", t: "14:32", me: false },
  { id: 2, text: "Привет! Да, давай. NTC или Tennis Club?", t: "14:35", me: true },
  { id: 3, text: "NTC, корт 3 свободен", t: "14:36", me: false },
  { id: 4, text: "Отлично, забронирую. 2 сета?", t: "14:38", me: true },
  { id: 5, text: "Да, best of 3 давай 💪", t: "14:39", me: false },
  { id: 6, text: "Договорились! До завтра", t: "14:40", me: true },
];

// ── Shared ──
function Av({ text, size = 40, color = C.primary, online }) {
  return (
    <div style={{ position: "relative", flexShrink: 0 }}>
      <div style={{ width: size, height: size, borderRadius: size/2, background: color+"20", display: "flex", alignItems: "center", justifyContent: "center", fontSize: size*0.35, fontWeight: 700, color }}>{text}</div>
      {online && <div style={{ position: "absolute", bottom: 0, right: 0, width: 12, height: 12, borderRadius: 6, background: C.green, border: "2px solid white" }}/>}
    </div>
  );
}

function Cd({ children, style, onClick }) {
  return <div onClick={onClick} style={{ background: C.card, borderRadius: 16, padding: 16, border: `1px solid ${C.borderL}`, cursor: onClick?"pointer":"default", ...style }}>{children}</div>;
}

function Seg({ options, active, onChange }) {
  return (
    <div style={{ display: "flex", background: C.borderL, borderRadius: 10, padding: 3 }}>
      {options.map(o => (
        <button key={o.id} onClick={() => onChange(o.id)} style={{ flex: 1, padding: "8px 10px", borderRadius: 8, border: "none", background: active===o.id?C.card:"transparent", color: active===o.id?C.text:C.sub, fontSize: 13, fontWeight: active===o.id?600:500, cursor: "pointer", boxShadow: active===o.id?"0 1px 3px rgba(0,0,0,0.08)":"none" }}>{o.label}</button>
      ))}
    </div>
  );
}

function Bdg({ text, color = C.primary, icon }) {
  return <span style={{ display: "inline-flex", alignItems: "center", gap: 4, padding: "3px 10px", borderRadius: 20, background: color+"15", color, fontSize: 12, fontWeight: 600 }}>{icon}{text}</span>;
}

function Chip({ label, active, onClick }) {
  return <button onClick={onClick} style={{ padding: "6px 14px", borderRadius: 20, border: `1px solid ${active?C.primary:C.border}`, background: active?C.primaryLight:C.card, color: active?C.primary:C.sub, fontSize: 13, fontWeight: 500, cursor: "pointer", whiteSpace: "nowrap" }}>{label}</button>;
}

// ── Phone Frame ──
function Frame({ children }) {
  return (
    <div style={{ width: 390, height: 844, background: C.bg, borderRadius: 40, overflow: "hidden", boxShadow: "0 25px 80px rgba(0,0,0,0.2)", position: "relative", fontFamily: "'SF Pro Display',-apple-system,BlinkMacSystemFont,sans-serif", border: "8px solid #1a1a1a" }}>
      <div style={{ position: "absolute", top: 0, left: "50%", transform: "translateX(-50%)", width: 126, height: 28, background: "#1a1a1a", borderRadius: "0 0 16px 16px", zIndex: 50 }}/>
      <div style={{ height: "100%", overflow: "hidden", display: "flex", flexDirection: "column" }}>{children}</div>
    </div>
  );
}

function SBar() {
  return <div style={{ height: 54, display: "flex", alignItems: "flex-end", justifyContent: "space-between", padding: "0 24px 6px", background: C.card, fontSize: 14, fontWeight: 600 }}><span>9:41</span><span style={{ fontSize: 13 }}>●●● 100% 🔋</span></div>;
}

function Header({ title, onBack, right }) {
  return (
    <div style={{ height: 52, display: "flex", alignItems: "center", justifyContent: "space-between", padding: "0 16px", background: C.card, borderBottom: `1px solid ${C.borderL}` }}>
      <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
        {onBack && <button onClick={onBack} style={{ background: "none", border: "none", fontSize: 22, cursor: "pointer", padding: 4 }}>←</button>}
        <span style={{ fontSize: 18, fontWeight: 700, color: C.text }}>{title}</span>
      </div>
      <div style={{ display: "flex", gap: 8, alignItems: "center" }}>{right}</div>
    </div>
  );
}

function TabBar({ active, onChange }) {
  const tabs = [["home","🏠","Главная"],["players","👥","Игроки"],["events","🎾","Ивенты"],["communities","🏛","Клубы"],["profile","👤","Профиль"]];
  return (
    <div style={{ height: 80, display: "flex", alignItems: "flex-start", justifyContent: "space-around", padding: "8px 0 0", background: C.card, borderTop: `1px solid ${C.borderL}` }}>
      {tabs.map(([id,ic,lb]) => (
        <button key={id} onClick={() => onChange(id)} style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: 2, background: "none", border: "none", cursor: "pointer", padding: "4px 8px" }}>
          <span style={{ fontSize: 22, opacity: active===id?1:0.4 }}>{ic}</span>
          <span style={{ fontSize: 10, fontWeight: active===id?700:500, color: active===id?C.primary:C.muted }}>{lb}</span>
        </button>
      ))}
    </div>
  );
}

function HeaderBadges({ navigate }) {
  return (
    <>
      <button onClick={() => navigate("chatList")} style={{ position: "relative", background: "none", border: "none", fontSize: 20, cursor: "pointer", padding: 4 }}>
        💬<span style={{ position: "absolute", top: 0, right: -2, width: 16, height: 16, borderRadius: 8, background: C.primary, color: "white", fontSize: 9, fontWeight: 700, display: "flex", alignItems: "center", justifyContent: "center" }}>2</span>
      </button>
      <button onClick={() => navigate("notifs")} style={{ position: "relative", background: "none", border: "none", fontSize: 20, cursor: "pointer", padding: 4 }}>
        🔔<span style={{ position: "absolute", top: 0, right: -2, width: 16, height: 16, borderRadius: 8, background: C.red, color: "white", fontSize: 9, fontWeight: 700, display: "flex", alignItems: "center", justifyContent: "center" }}>5</span>
      </button>
    </>
  );
}

// ── AUTH ──
function AuthPhone({ onNext }) {
  const [ph, setPh] = useState("");
  return (
    <div style={{ flex: 1, display: "flex", flexDirection: "column", padding: 24, paddingTop: 60, background: C.card }}>
      <div style={{ fontSize: 48, marginBottom: 12 }}>🎾</div>
      <h1 style={{ fontSize: 28, fontWeight: 800, margin: 0 }}>Tennis Astana</h1>
      <p style={{ fontSize: 15, color: C.sub, margin: "8px 0 40px" }}>Найди партнёра. Играй. Расти.</p>
      <label style={{ fontSize: 13, fontWeight: 600, color: C.sub, marginBottom: 8 }}>Номер телефона</label>
      <div style={{ display: "flex", alignItems: "center", gap: 10, padding: "14px 16px", borderRadius: 14, border: `2px solid ${ph?C.primary:C.border}`, background: C.bg }}>
        <span style={{ fontSize: 16, fontWeight: 600 }}>🇰🇿 +7</span>
        <input value={ph} onChange={e=>setPh(e.target.value.replace(/\D/g,"").slice(0,10))} placeholder="707 123 45 67" style={{ flex: 1, border: "none", outline: "none", fontSize: 18, fontWeight: 500, background: "transparent", letterSpacing: 1 }}/>
      </div>
      <p style={{ fontSize: 12, color: C.muted, margin: "12px 0" }}>Мы отправим SMS с кодом подтверждения</p>
      <div style={{ flex: 1 }}/>
      <button onClick={() => ph.length>=10&&onNext()} style={{ width: "100%", padding: "16px", borderRadius: 14, border: "none", background: ph.length>=10?C.primary:C.border, color: ph.length>=10?"white":C.muted, fontSize: 17, fontWeight: 700, cursor: ph.length>=10?"pointer":"default" }}>Получить код</button>
      <div style={{ height: 34 }}/>
    </div>
  );
}

function AuthOTP({ onNext, onBack }) {
  const [code, setCode] = useState(["","","",""]);
  const [timer, setTimer] = useState(45);
  useEffect(() => { const t = setInterval(() => setTimer(p => p>0?p-1:0), 1000); return () => clearInterval(t); }, []);
  return (
    <div style={{ flex: 1, display: "flex", flexDirection: "column", padding: 24, paddingTop: 20, background: C.card }}>
      <button onClick={onBack} style={{ background: "none", border: "none", fontSize: 18, cursor: "pointer", alignSelf: "flex-start", padding: "8px 0" }}>← Назад</button>
      <h1 style={{ fontSize: 26, fontWeight: 800, margin: "24px 0 8px" }}>Введите код</h1>
      <p style={{ fontSize: 15, color: C.sub }}>Отправлен на +7 707 *** 45 67</p>
      <div style={{ display: "flex", gap: 12, justifyContent: "center", margin: "40px 0" }}>
        {code.map((c,i) => (
          <div key={i} style={{ width: 64, height: 72, borderRadius: 16, border: `2px solid ${c?C.primary:C.border}`, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 28, fontWeight: 700, background: c?C.primaryLight:C.bg }}>{c||"•"}</div>
        ))}
      </div>
      <div style={{ textAlign: "center" }}>{timer>0?<span style={{ color: C.muted, fontSize: 14 }}>Повторно через {timer}с</span>:<button style={{ background: "none", border: "none", color: C.primary, fontWeight: 600, cursor: "pointer" }}>Отправить ещё раз</button>}</div>
      <div style={{ flex: 1 }}/>
      <div style={{ display: "grid", gridTemplateColumns: "repeat(3,1fr)", gap: 8, padding: "0 16px" }}>
        {[1,2,3,4,5,6,7,8,9,null,0,"⌫"].map((n,i) => (
          <button key={i} onClick={() => {
            if(n===null)return;
            if(n==="⌫"){const nc=[...code];const li=code.findLastIndex(c=>c!=="");if(li>=0)nc[li]="";setCode(nc);}
            else{const fi=code.findIndex(c=>c==="");if(fi>=0){const nc=[...code];nc[fi]=String(n);setCode(nc);if(fi===3)setTimeout(onNext,300);}}
          }} style={{ height: 52, borderRadius: 12, border: "none", background: n===null?"transparent":C.borderL, fontSize: 22, fontWeight: 600, cursor: n===null?"default":"pointer" }}>{n}</button>
        ))}
      </div>
      <div style={{ height: 20 }}/>
    </div>
  );
}

function AuthProfile({ onNext }) {
  const [step, setStep] = useState(0);
  if(step===0) return (
    <div style={{ flex: 1, display: "flex", flexDirection: "column", padding: 24, paddingTop: 40, background: C.card }}>
      <h1 style={{ fontSize: 26, fontWeight: 800, margin: 0 }}>Расскажите о себе</h1>
      <p style={{ fontSize: 15, color: C.sub, margin: "8px 0 32px" }}>Чтобы находить подходящих партнёров</p>
      {[["Имя","Иван"],["Фамилия","Петров"]].map(([l,p])=>(
        <div key={l} style={{ marginBottom: 20 }}>
          <label style={{ fontSize: 13, fontWeight: 600, color: C.sub }}>{l}</label>
          <input placeholder={p} style={{ width: "100%", padding: "14px 16px", borderRadius: 12, border: `1.5px solid ${C.border}`, fontSize: 16, marginTop: 6, outline: "none", boxSizing: "border-box" }}/>
        </div>
      ))}
      <label style={{ fontSize: 13, fontWeight: 600, color: C.sub }}>Пол</label>
      <div style={{ display: "flex", gap: 10, margin: "8px 0 20px" }}>
        {["👨 Мужской","👩 Женский"].map(g=><button key={g} style={{ flex: 1, padding: 14, borderRadius: 12, border: `1.5px solid ${C.border}`, background: C.bg, fontSize: 15, cursor: "pointer" }}>{g}</button>)}
      </div>
      <div style={{ flex: 1 }}/>
      <button onClick={()=>setStep(1)} style={{ width: "100%", padding: 16, borderRadius: 14, border: "none", background: C.primary, color: "white", fontSize: 17, fontWeight: 700, cursor: "pointer" }}>Далее</button>
      <div style={{ height: 34 }}/>
    </div>
  );
  return (
    <div style={{ flex: 1, display: "flex", flexDirection: "column", padding: 24, paddingTop: 40, background: C.card }}>
      <div style={{ display: "flex", gap: 4, marginBottom: 32 }}>{[0,1,2].map(i=><div key={i} style={{ flex: 1, height: 4, borderRadius: 2, background: i===0?C.primary:C.borderL }}/>)}</div>
      <p style={{ fontSize: 13, color: C.primary, fontWeight: 600, margin: 0 }}>Вопрос 1 из 3</p>
      <h2 style={{ fontSize: 22, fontWeight: 700, margin: "8px 0 24px" }}>Как давно вы играете?</h2>
      {[{e:"🌱",t:"Менее 1 года",s:"Только начинаю"},{e:"🌿",t:"1-3 года",s:"Базовые навыки"},{e:"🌳",t:"3-5 лет",s:"Уверенная игра"},{e:"🏆",t:"Более 5 лет",s:"Опытный игрок"}].map((o,i)=>(
        <button key={i} onClick={()=>setTimeout(onNext,300)} style={{ display: "flex", alignItems: "center", gap: 14, padding: 16, borderRadius: 14, border: `1.5px solid ${C.border}`, background: C.bg, marginBottom: 10, cursor: "pointer", textAlign: "left" }}>
          <span style={{ fontSize: 28 }}>{o.e}</span>
          <div><div style={{ fontSize: 16, fontWeight: 600 }}>{o.t}</div><div style={{ fontSize: 13, color: C.muted }}>{o.s}</div></div>
        </button>
      ))}
    </div>
  );
}

// ── HOME ──
function HomeScreen({ nav }) {
  const [ft, setFt] = useState("news");
  return (
    <div style={{ flex: 1, overflow: "auto" }}>
      <div style={{ padding: 16, display: "flex", flexDirection: "column", gap: 16 }}>
        <Cd style={{ background: `linear-gradient(135deg,${C.primary},${C.primaryDark})`, border: "none", color: "white" }}>
          <div style={{ display: "flex", justifyContent: "space-between" }}>
            <div>
              <p style={{ fontSize: 13, opacity: 0.8, margin: 0 }}>Мой рейтинг</p>
              <p style={{ fontSize: 36, fontWeight: 800, margin: "4px 0" }}>1,250</p>
              <div style={{ display: "flex", alignItems: "center", gap: 6 }}>
                <span style={{ background: "rgba(255,255,255,0.2)", padding: "3px 10px", borderRadius: 20, fontSize: 12, fontWeight: 600 }}>🎾 Любитель 3.0</span>
                <span style={{ fontSize: 13, color: "#A7F3D0" }}>↑ 18.5</span>
              </div>
            </div>
            <div style={{ textAlign: "right" }}>
              <p style={{ fontSize: 13, opacity: 0.8, margin: 0 }}>Позиция</p>
              <p style={{ fontSize: 24, fontWeight: 700, margin: "4px 0" }}>#45</p>
              <p style={{ fontSize: 12, opacity: 0.7, margin: 0 }}>из 500</p>
            </div>
          </div>
        </Cd>
        <div style={{ display: "flex", gap: 10 }}>
          <button onClick={()=>nav("eventCreate")} style={{ flex: 1, padding: "14px", borderRadius: 14, border: `2px solid ${C.primary}`, background: C.primaryLight, color: C.primary, fontSize: 14, fontWeight: 700, cursor: "pointer" }}>🎾 Найти партнёра</button>
          <button onClick={()=>nav("eventCreate")} style={{ flex: 1, padding: "14px", borderRadius: 14, border: `2px solid ${C.accent}`, background: C.accent+"15", color: "#B8860B", fontSize: 14, fontWeight: 700, cursor: "pointer" }}>➕ Создать игру</button>
        </div>
        <div>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 10 }}>
            <span style={{ fontSize: 16, fontWeight: 700 }}>Ближайшие игры</span>
            <button onClick={()=>nav("events")} style={{ background: "none", border: "none", color: C.primary, fontWeight: 600, fontSize: 13, cursor: "pointer" }}>Все →</button>
          </div>
          <Cd onClick={()=>nav("eventDetail")} style={{ display: "flex", alignItems: "center", gap: 12 }}>
            <div style={{ width: 48, height: 48, borderRadius: 12, background: C.primaryLight, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 22 }}>🎾</div>
            <div style={{ flex: 1 }}><p style={{ fontSize: 15, fontWeight: 600, margin: 0 }}>Парная игра выходного дня</p><p style={{ fontSize: 13, color: C.muted, margin: "2px 0 0" }}>Сб, 10:00 · Tennis Club</p></div>
            <Bdg text="6/8" color={C.yellow}/>
          </Cd>
        </div>
        <Seg options={[{id:"news",label:"Новости"},{id:"feed",label:"Лента"}]} active={ft} onChange={setFt}/>
        <Cd>
          <div style={{ display: "flex", alignItems: "center", gap: 10, marginBottom: 10 }}>
            <Av text="NT"/><div style={{ flex: 1 }}><span style={{ fontWeight: 600, fontSize: 14 }}>NTC Astana</span><span style={{ fontSize: 12, color: C.muted, marginLeft: 8 }}>2ч</span></div>
            <span style={{ fontSize: 11, color: C.primary, fontWeight: 600, background: C.primaryLight, padding: "2px 8px", borderRadius: 8 }}>✓ Клуб</span>
          </div>
          <p style={{ fontSize: 14, margin: 0, lineHeight: 1.5 }}>🏆 Открыта регистрация на весенний турнир! 16 участников, призовой фонд — сертификаты.</p>
          <div style={{ display: "flex", gap: 16, marginTop: 12, paddingTop: 10, borderTop: `1px solid ${C.borderL}` }}><span style={{ fontSize: 13, color: C.muted }}>❤️ 24</span><span style={{ fontSize: 13, color: C.muted }}>💬 8</span></div>
        </Cd>
        <Cd style={{ background: `linear-gradient(135deg,${C.primaryLight},#f0fdf4)`, border: `1px solid ${C.primary}22` }}>
          <div style={{ display: "flex", alignItems: "center", gap: 8, marginBottom: 8 }}><span style={{ fontSize: 16 }}>🎾</span><span style={{ fontWeight: 700, fontSize: 14, color: C.primary }}>Результат матча</span></div>
          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
            <div style={{ textAlign: "center" }}><Av text="ИП" size={36}/><p style={{ fontSize: 12, fontWeight: 600, margin: "4px 0 0" }}>Иван П.</p></div>
            <div style={{ textAlign: "center" }}><span style={{ fontSize: 24, fontWeight: 800, color: C.primary }}>6-4 3-6 7-5</span><p style={{ fontSize: 11, color: C.muted, margin: "4px 0 0" }}>Вчера · NTC</p></div>
            <div style={{ textAlign: "center" }}><Av text="АМ" size={36} color={C.sub}/><p style={{ fontSize: 12, fontWeight: 600, margin: "4px 0 0" }}>Алексей М.</p></div>
          </div>
        </Cd>
      </div>
    </div>
  );
}

// ── PLAYERS ──
function PlayersScreen({ nav }) {
  return (
    <div style={{ flex: 1, overflow: "auto" }}>
      <div style={{ padding: "12px 16px", position: "sticky", top: 0, background: C.bg, zIndex: 10 }}>
        <div style={{ display: "flex", alignItems: "center", gap: 8, padding: "10px 14px", borderRadius: 12, background: C.card, border: `1px solid ${C.border}` }}>
          <span style={{ color: C.muted }}>🔍</span>
          <input placeholder="Поиск игроков..." style={{ flex: 1, border: "none", outline: "none", fontSize: 15, background: "transparent" }}/>
        </div>
        <div style={{ display: "flex", gap: 8, marginTop: 10, overflowX: "auto" }}>
          {["Все","3.0-4.0","Есильский","Мужской","Онлайн"].map((f,i)=><Chip key={f} label={f} active={i===0}/>)}
        </div>
      </div>
      <div style={{ padding: "0 16px 16px", display: "flex", flexDirection: "column", gap: 8 }}>
        {PLAYERS.map(p=>(
          <Cd key={p.id} onClick={()=>nav("playerProfile")} style={{ display: "flex", alignItems: "center", gap: 12 }}>
            <Av text={p.av} size={48} color={p.g==="Ж"?"#EC4899":C.primary} online={p.on}/>
            <div style={{ flex: 1 }}>
              <span style={{ fontWeight: 600, fontSize: 15 }}>{p.name}</span>
              <div style={{ display: "flex", gap: 8, marginTop: 4 }}>
                <span style={{ fontSize: 12, color: C.muted }}>🎾 {p.lv}</span>
                <span style={{ fontSize: 12, color: C.muted }}>📊 {p.r}</span>
                <span style={{ fontSize: 12, color: C.muted }}>📍 {p.dist}</span>
              </div>
            </div>
            <div style={{ textAlign: "right" }}><span style={{ fontSize: 14, fontWeight: 700, color: C.primary }}>{p.wr}%</span><p style={{ fontSize: 11, color: C.muted, margin: "2px 0 0" }}>{p.games} игр</p></div>
          </Cd>
        ))}
      </div>
    </div>
  );
}

// ── EVENTS ──
function EventsScreen({ nav }) {
  const [tab, setTab] = useState("feed");
  const [fl, setFl] = useState("all");
  return (
    <div style={{ flex: 1, overflow: "auto", position: "relative" }}>
      <div style={{ padding: "12px 16px 0", background: C.bg, position: "sticky", top: 0, zIndex: 10 }}>
        <Seg options={[{id:"feed",label:"Лента"},{id:"calendar",label:"Календарь"},{id:"my",label:"Мои"}]} active={tab} onChange={setTab}/>
        {tab==="feed"&&<div style={{ display: "flex", gap: 8, marginTop: 10, paddingBottom: 10, overflowX: "auto" }}>
          {[["all","Все"],["find_partner","Партнёр"],["organized_game","Игра"],["tournament","Турнир"],["training","Тренировка"]].map(([id,lb])=><Chip key={id} label={lb} active={fl===id} onClick={()=>setFl(id)}/>)}
        </div>}
      </div>
      {tab==="feed"&&<div style={{ padding: "8px 16px 16px", display: "flex", flexDirection: "column", gap: 10 }}>
        {EVENTS.filter(e=>fl==="all"||e.type===fl).map(e=>(
          <Cd key={e.id} onClick={()=>nav("eventDetail")} style={{ padding: 14 }}>
            <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 8 }}>
              <div style={{ flex: 1 }}>
                <div style={{ display: "flex", alignItems: "center", gap: 6, marginBottom: 4 }}>
                  <div style={{ width: 8, height: 8, borderRadius: 4, background: e.sc }}/>
                  <span style={{ fontSize: 12, color: C.muted }}>{e.comp} · {e.level}</span>
                </div>
                <p style={{ fontSize: 15, fontWeight: 700, margin: 0 }}>{e.title}</p>
              </div>
              <span style={{ fontSize: 15, fontWeight: 700, color: C.primary }}>{e.spots}</span>
            </div>
            <div style={{ display: "flex", gap: 12, fontSize: 13, color: C.sub }}><span>🕐 {e.time}</span><span>📍 {e.loc}</span></div>
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginTop: 8, paddingTop: 8, borderTop: `1px solid ${C.borderL}` }}>
              <span style={{ fontSize: 12, color: C.muted }}>{e.creator}</span>
              <button style={{ padding: "6px 16px", borderRadius: 8, border: "none", background: C.primary, color: "white", fontSize: 13, fontWeight: 600, cursor: "pointer" }}>Записаться</button>
            </div>
          </Cd>
        ))}
      </div>}
      {tab==="calendar"&&<div style={{ padding: 16 }}>
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 16 }}>
          <span style={{ fontSize: 18, fontWeight: 700 }}>Март 2026</span>
          <div style={{ display: "flex", gap: 8 }}><button style={{ width: 32, height: 32, borderRadius: 8, border: `1px solid ${C.border}`, background: C.card, cursor: "pointer" }}>←</button><button style={{ width: 32, height: 32, borderRadius: 8, border: `1px solid ${C.border}`, background: C.card, cursor: "pointer" }}>→</button></div>
        </div>
        <div style={{ display: "grid", gridTemplateColumns: "repeat(7,1fr)", gap: 4, textAlign: "center" }}>
          {["Пн","Вт","Ср","Чт","Пт","Сб","Вс"].map(d=><span key={d} style={{ fontSize: 11, color: C.muted, padding: 4 }}>{d}</span>)}
          {Array.from({length:31},(_,i)=>i+1).map(d=>{
            const has=[3,8,10,15,16,22,28].includes(d), today=d===10;
            return <div key={d} style={{ padding: 6, borderRadius: 8, background: today?C.primary:"transparent", color: today?"white":C.text, fontWeight: today?700:400, fontSize: 14, cursor: "pointer", position: "relative" }}>{d}{has&&<div style={{ width: 4, height: 4, borderRadius: 2, background: today?"white":C.primary, margin: "2px auto 0" }}/>}</div>;
          })}
        </div>
        <div style={{ marginTop: 20 }}>
          <p style={{ fontSize: 15, fontWeight: 700, margin: "0 0 10px" }}>10 марта</p>
          <Cd onClick={()=>nav("eventDetail")} style={{ display: "flex", alignItems: "center", gap: 10 }}>
            <div style={{ width: 4, height: 36, borderRadius: 2, background: C.primary }}/>
            <div><p style={{ fontSize: 14, fontWeight: 600, margin: 0 }}>Ищу партнёра на вечер</p><p style={{ fontSize: 12, color: C.muted, margin: "2px 0 0" }}>18:00 · NTC Astana</p></div>
          </Cd>
        </div>
      </div>}
      {tab==="my"&&<div style={{ padding: 16, textAlign: "center", color: C.muted, paddingTop: 60 }}>
        <span style={{ fontSize: 48 }}>🎾</span><p style={{ fontSize: 16, fontWeight: 600, color: C.text, margin: "16px 0 4px" }}>Нет записей</p><p style={{ fontSize: 14 }}>Запишитесь на ивент или создайте свой</p>
        <button onClick={()=>nav("eventCreate")} style={{ marginTop: 16, padding: "12px 24px", borderRadius: 12, border: "none", background: C.primary, color: "white", fontWeight: 600, cursor: "pointer" }}>Создать ивент</button>
      </div>}
      <button onClick={()=>nav("eventCreate")} style={{ position: "absolute", bottom: 16, right: 20, width: 56, height: 56, borderRadius: 28, background: C.primary, border: "none", color: "white", fontSize: 28, cursor: "pointer", boxShadow: "0 4px 16px rgba(15,139,79,0.4)", display: "flex", alignItems: "center", justifyContent: "center", zIndex: 20 }}>+</button>
    </div>
  );
}

// ── EVENT DETAIL ──
function EventDetail({ onBack }) {
  return (
    <div style={{ flex: 1, overflow: "auto" }}>
      <div style={{ padding: 16, display: "flex", flexDirection: "column", gap: 14 }}>
        <div style={{ display: "flex", gap: 8 }}><Bdg icon="🟡" text="Идёт набор" color={C.yellow}/><Bdg text="Doubles" color={C.blue}/></div>
        <h2 style={{ fontSize: 22, fontWeight: 800, margin: 0 }}>Парная игра выходного дня</h2>
        <div style={{ display: "flex", flexWrap: "wrap", gap: 8, fontSize: 14, color: C.sub }}><span>📅 Сб, 10 марта, 10:00–12:00</span><span>📍 Tennis Club Astana</span><span>🎾 2.5–3.5</span></div>
        <Cd style={{ background: C.primaryLight, border: `1px solid ${C.primary}22` }}>
          <div style={{ display: "flex", justifyContent: "space-around", textAlign: "center" }}>
            <div><p style={{ fontSize: 24, fontWeight: 800, color: C.primary, margin: 0 }}>6/8</p><p style={{ fontSize: 12, color: C.sub, margin: "2px 0 0" }}>Участников</p></div>
            <div style={{ width: 1, background: C.primary+"30" }}/>
            <div><p style={{ fontSize: 24, fontWeight: 800, color: C.primary, margin: 0 }}>2</p><p style={{ fontSize: 12, color: C.sub, margin: "2px 0 0" }}>Места</p></div>
            <div style={{ width: 1, background: C.primary+"30" }}/>
            <div><p style={{ fontSize: 24, fontWeight: 800, color: C.primary, margin: 0 }}>3</p><p style={{ fontSize: 12, color: C.sub, margin: "2px 0 0" }}>Сета</p></div>
          </div>
        </Cd>
        <p style={{ fontSize: 14, lineHeight: 1.6, margin: 0 }}>Приглашаем на парную игру! Формат: 3 сета, тай-брейк в решающем. Мячи предоставляются.</p>
        <div>
          <span style={{ fontWeight: 700, fontSize: 15 }}>Участники (6)</span>
          <div style={{ display: "flex", flexWrap: "wrap", gap: 8, marginTop: 10 }}>
            {PLAYERS.map(p=><div key={p.id} style={{ display: "flex", alignItems: "center", gap: 6, padding: "6px 12px", borderRadius: 20, background: C.borderL }}><Av text={p.av} size={24}/><span style={{ fontSize: 13, fontWeight: 500 }}>{p.name.split(" ")[0]}</span></div>)}
            <div style={{ padding: "6px 12px", borderRadius: 20, background: C.borderL, border: `1px dashed ${C.border}` }}><span style={{ fontSize: 13, color: C.muted }}>+2 места</span></div>
          </div>
        </div>
        <Cd style={{ display: "flex", alignItems: "center", gap: 12 }}>
          <Av text="NT"/><div><div style={{ display: "flex", alignItems: "center", gap: 6 }}><span style={{ fontWeight: 600 }}>NTC Astana</span><span style={{ color: C.primary }}>✓</span></div><span style={{ fontSize: 12, color: C.muted }}>Клуб · 245 участников</span></div>
        </Cd>
      </div>
      <div style={{ padding: "12px 16px 32px", position: "sticky", bottom: 0, background: C.bg, borderTop: `1px solid ${C.borderL}` }}>
        <button style={{ width: "100%", padding: 16, borderRadius: 14, border: "none", background: C.primary, color: "white", fontSize: 17, fontWeight: 700, cursor: "pointer" }}>Записаться</button>
      </div>
    </div>
  );
}

// ── EVENT CREATE ──
function EventCreate({ onBack }) {
  const [step, setStep] = useState(0);
  const steps = ["Тип","Состав","Формат","Место","Время","Участники","Превью"];
  return (
    <div style={{ flex: 1, overflow: "auto" }}>
      <div style={{ padding: "12px 16px", display: "flex", gap: 4 }}>{steps.map((_,i)=><div key={i} style={{ flex: 1, height: 4, borderRadius: 2, background: i<=step?C.primary:C.borderL }}/>)}</div>
      <div style={{ padding: "8px 16px 0" }}><p style={{ fontSize: 12, color: C.primary, fontWeight: 600, margin: 0 }}>Шаг {step+1} из {steps.length}</p><h2 style={{ fontSize: 22, fontWeight: 700, margin: "4px 0 20px" }}>{steps[step]}</h2></div>
      {step===0&&<div style={{ padding: "0 16px", display: "flex", flexDirection: "column", gap: 10 }}>
        {[{i:"🤝",t:"Ищу партнёра",d:"Найти соперника",c:C.primary},{i:"🎾",t:"Организую игру",d:"Групповая игра",c:C.blue},{i:"🏆",t:"Турнир",d:"Соревнование с сеткой",c:C.accent},{i:"📚",t:"Тренировка",d:"Групповое занятие",c:"#8B5CF6"}].map((t,i)=>(
          <button key={i} onClick={()=>setStep(1)} style={{ display: "flex", alignItems: "center", gap: 14, padding: 16, borderRadius: 14, border: `2px solid ${C.border}`, background: C.card, cursor: "pointer", textAlign: "left" }}>
            <div style={{ width: 48, height: 48, borderRadius: 14, background: t.c+"15", display: "flex", alignItems: "center", justifyContent: "center", fontSize: 24 }}>{t.i}</div>
            <div><p style={{ fontSize: 16, fontWeight: 700, margin: 0 }}>{t.t}</p><p style={{ fontSize: 13, color: C.muted, margin: "2px 0 0" }}>{t.d}</p></div>
          </button>
        ))}
      </div>}
      {step===1&&<div style={{ padding: "0 16px", display: "grid", gridTemplateColumns: "1fr 1fr", gap: 10 }}>
        {[{i:"👤",t:"Singles",d:"1 на 1"},{i:"👥",t:"Doubles",d:"2 на 2"},{i:"👫",t:"Mixed",d:"Смешанные"},{i:"🏟",t:"Team",d:"Командная"}].map((t,i)=>(
          <button key={i} onClick={()=>setStep(2)} style={{ padding: 20, borderRadius: 14, border: `2px solid ${C.border}`, background: C.card, cursor: "pointer", textAlign: "center" }}>
            <span style={{ fontSize: 32 }}>{t.i}</span><p style={{ fontSize: 15, fontWeight: 700, margin: "8px 0 2px" }}>{t.t}</p><p style={{ fontSize: 12, color: C.muted, margin: 0 }}>{t.d}</p>
          </button>
        ))}
      </div>}
      {step>=2&&step<6&&<div style={{ padding: "0 16px" }}><Cd style={{ padding: 20, textAlign: "center" }}>
        <span style={{ fontSize: 40 }}>{["⚙️","📍","🕐","👥"][step-2]}</span>
        <p style={{ fontSize: 15, fontWeight: 600, color: C.text, margin: "12px 0 4px" }}>{["Настройки формата","Выбор места","Дата и время","Ограничения"][step-2]}</p>
      </Cd></div>}
      {step===6&&<div style={{ padding: "0 16px" }}><Cd style={{ background: C.primaryLight, border: `1px solid ${C.primary}22` }}>
        <p style={{ fontSize: 11, color: C.primary, fontWeight: 600, margin: "0 0 8px" }}>ПРЕВЬЮ</p>
        <p style={{ fontSize: 18, fontWeight: 700, margin: "0 0 8px" }}>Ищу партнёра на вечер</p>
        <div style={{ display: "flex", flexDirection: "column", gap: 6, fontSize: 14 }}><span>🎾 Singles · 2 сета</span><span>📍 NTC Astana</span><span>📅 Сегодня, 18:00–20:00</span><span>🎯 Уровень 3.0–4.0</span></div>
      </Cd></div>}
      <div style={{ padding: "20px 16px 32px", display: "flex", gap: 10, position: "sticky", bottom: 0, background: C.bg }}>
        {step>0&&<button onClick={()=>setStep(step-1)} style={{ flex: 1, padding: 14, borderRadius: 14, border: `2px solid ${C.border}`, background: C.card, fontSize: 15, fontWeight: 600, cursor: "pointer" }}>Назад</button>}
        <button onClick={()=>step<6?setStep(step+1):onBack()} style={{ flex: 2, padding: 14, borderRadius: 14, border: "none", background: C.primary, color: "white", fontSize: 15, fontWeight: 700, cursor: "pointer" }}>{step===6?"Опубликовать":"Далее"}</button>
      </div>
    </div>
  );
}

// ── COMMUNITY ──
function CommunityScreen({ onBack }) {
  const [tab, setTab] = useState("feed");
  return (
    <div style={{ flex: 1, overflow: "auto" }}>
      <div style={{ height: 130, background: `linear-gradient(135deg,${C.primary},${C.primaryDark})`, display: "flex", alignItems: "flex-end", padding: "0 16px 16px" }}>
        <div style={{ display: "flex", alignItems: "center", gap: 12 }}>
          <div style={{ width: 52, height: 52, borderRadius: 14, background: "white", display: "flex", alignItems: "center", justifyContent: "center", fontSize: 22, fontWeight: 700, color: C.primary }}>NT</div>
          <div style={{ color: "white" }}><div style={{ display: "flex", alignItems: "center", gap: 6 }}><span style={{ fontSize: 20, fontWeight: 800 }}>NTC Astana</span><span>✓</span></div><span style={{ fontSize: 13, opacity: 0.8 }}>Клуб · 245 участников</span></div>
        </div>
      </div>
      <div style={{ display: "flex", borderBottom: `2px solid ${C.borderL}`, padding: "0 8px", background: C.card, overflowX: "auto" }}>
        {[["feed","Лента"],["events","Ивенты"],["rating","Рейтинг"],["members","Участники"],["chat","Чат"],["photos","Фото"]].map(([id,lb])=>(
          <button key={id} onClick={()=>setTab(id)} style={{ padding: "12px 12px", border: "none", borderBottom: `2px solid ${tab===id?C.primary:"transparent"}`, background: "none", color: tab===id?C.primary:C.sub, fontWeight: tab===id?700:500, fontSize: 13, cursor: "pointer", whiteSpace: "nowrap" }}>{lb}</button>
        ))}
      </div>
      <div style={{ padding: 16 }}>
        {tab==="feed"&&<Cd><div style={{ display: "flex", alignItems: "center", gap: 10, marginBottom: 8 }}><Av text="NT" size={32}/><span style={{ fontWeight: 600, fontSize: 14 }}>NTC Astana</span><span style={{ fontSize: 12, color: C.muted }}>5ч</span></div><p style={{ fontSize: 14, margin: 0, lineHeight: 1.5 }}>📢 С марта корты работают до 23:00!</p><div style={{ display: "flex", gap: 16, marginTop: 10, paddingTop: 8, borderTop: `1px solid ${C.borderL}`, fontSize: 13, color: C.muted }}><span>❤️ 18</span><span>💬 5</span></div></Cd>}
        {tab==="rating"&&<div style={{ display: "flex", flexDirection: "column", gap: 6 }}>
          {[{r:1,n:"Марат К.",rt:1450,g:32,w:75},{r:2,n:"Дмитрий Л.",rt:1380,g:28,w:68},{r:3,n:"Алексей М.",rt:1320,g:45,w:62},{r:4,n:"Иван П.",rt:1250,g:20,w:60},{r:5,n:"Арман Ж.",rt:1200,g:15,w:53}].map(p=>(
            <div key={p.r} style={{ display: "flex", alignItems: "center", gap: 12, padding: "10px 12px", borderRadius: 12, background: p.r<=3?C.primaryLight:C.card }}>
              <span style={{ width: 28, textAlign: "center", fontSize: p.r<=3?20:16, fontWeight: 700, color: p.r<=3?C.primary:C.muted }}>{p.r<=3?["🥇","🥈","🥉"][p.r-1]:p.r}</span>
              <Av text={p.n.split(" ").map(w=>w[0]).join("")} size={36}/><div style={{ flex: 1 }}><span style={{ fontWeight: 600, fontSize: 14 }}>{p.n}</span><div style={{ fontSize: 12, color: C.muted }}>{p.g} игр · {p.w}%</div></div>
              <span style={{ fontSize: 16, fontWeight: 800, color: C.primary }}>{p.rt}</span>
            </div>
          ))}
        </div>}
        {tab==="members"&&<div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
          {PLAYERS.slice(0,4).map(p=><div key={p.id} style={{ display: "flex", alignItems: "center", gap: 10, padding: 8 }}><Av text={p.av} size={40} online={p.on}/><div style={{ flex: 1 }}><span style={{ fontWeight: 600, fontSize: 14 }}>{p.name}</span><span style={{ fontSize: 12, color: C.muted, marginLeft: 8 }}>🎾 {p.lv}</span></div><Bdg text="Участник" color={C.muted}/></div>)}
        </div>}
        {["events","chat","photos"].includes(tab)&&<div style={{ textAlign: "center", padding: "40px 0", color: C.muted }}><span style={{ fontSize: 40 }}>{tab==="events"?"🎾":tab==="chat"?"💬":"📷"}</span><p style={{ fontSize: 15, fontWeight: 600, color: C.text, margin: "12px 0 4px" }}>{tab==="events"?"Ивенты клуба":tab==="chat"?"Групповой чат":"Фотогалерея"}</p></div>}
      </div>
    </div>
  );
}

// ── PROFILE ──
function ProfileScreen({ nav }) {
  return (
    <div style={{ flex: 1, overflow: "auto" }}>
      <div style={{ padding: "20px 16px", background: C.card, textAlign: "center" }}>
        <div style={{ width: 80, height: 80, borderRadius: 40, background: C.primaryLight, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 32, fontWeight: 700, color: C.primary, margin: "0 auto 12px", border: `3px solid ${C.primary}` }}>ИП</div>
        <h2 style={{ fontSize: 22, fontWeight: 800, margin: 0 }}>Иван Петров</h2>
        <p style={{ fontSize: 14, color: C.sub, margin: "4px 0" }}>🎾 Любитель 3.0 · Есильский</p>
        <div style={{ display: "flex", justifyContent: "center", gap: 24, marginTop: 16 }}>
          {[["34","Игр",C.primary],["22","Побед",C.green],["64.7%","Win Rate",C.text],["3","🔥 Серия",C.accent]].map(([v,l,c])=><div key={l} style={{ textAlign: "center" }}><p style={{ fontSize: 22, fontWeight: 800, color: c, margin: 0 }}>{v}</p><p style={{ fontSize: 12, color: C.muted, margin: 0 }}>{l}</p></div>)}
        </div>
      </div>
      <div style={{ padding: 16 }}>
        <Cd>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 12 }}><span style={{ fontWeight: 700, fontSize: 15 }}>Рейтинг</span><span style={{ fontSize: 24, fontWeight: 800, color: C.primary }}>1,250</span></div>
          <div style={{ display: "flex", alignItems: "flex-end", height: 80, gap: 3, paddingBottom: 4 }}>
            {[40,45,42,50,48,55,52,58,62,60,65,70,68,72,75,78,74,80,82,85].map((h,i)=><div key={i} style={{ flex: 1, height: `${h}%`, background: i===19?C.primary:C.primary+"40", borderRadius: 2 }}/>)}
          </div>
          <div style={{ display: "flex", justifyContent: "space-between", fontSize: 11, color: C.muted, marginTop: 4 }}><span>Окт</span><span>Ноя</span><span>Дек</span><span>Янв</span><span>Фев</span><span>Мар</span></div>
          <div style={{ display: "flex", alignItems: "center", gap: 4, marginTop: 8, justifyContent: "center" }}><span style={{ fontSize: 14, color: C.green, fontWeight: 600 }}>↑ 250</span><span style={{ fontSize: 12, color: C.muted }}>за 6 мес.</span></div>
        </Cd>
      </div>
      <div style={{ padding: "0 16px" }}>
        <span style={{ fontWeight: 700, fontSize: 15 }}>Достижения ({BADGES.filter(b=>b.earned).length}/{BADGES.length})</span>
        <div style={{ display: "grid", gridTemplateColumns: "repeat(5,1fr)", gap: 8, marginTop: 12 }}>
          {BADGES.map(b=>(
            <div key={b.id} style={{ textAlign: "center", opacity: b.earned?1:0.4, position: "relative" }}>
              <div style={{ width: 48, height: 48, borderRadius: 14, background: b.earned?C.primaryLight:C.borderL, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 24, margin: "0 auto" }}>{b.icon}</div>
              <p style={{ fontSize: 10, margin: "4px 0 0", fontWeight: 500 }}>{b.name}</p>
              {!b.earned&&<div style={{ position: "absolute", bottom: 18, left: "50%", transform: "translateX(-50%)", width: 32, height: 3, borderRadius: 2, background: C.borderL }}><div style={{ width: `${b.p}%`, height: "100%", background: C.primary, borderRadius: 2 }}/></div>}
            </div>
          ))}
        </div>
      </div>
      <div style={{ padding: "20px 16px", display: "flex", flexDirection: "column", gap: 4 }}>
        {[["📊","История матчей","34"],["🏛","Мои сообщества","3"],["👥","Друзья","12"],["📝","Мои посты","5"],["⚙️","Настройки",""]].map(([ic,lb,ct])=>(
          <button key={lb} style={{ display: "flex", alignItems: "center", gap: 12, padding: "14px 12px", borderRadius: 12, border: "none", background: C.card, cursor: "pointer", width: "100%", textAlign: "left" }}>
            <span style={{ fontSize: 20 }}>{ic}</span><span style={{ flex: 1, fontSize: 15, fontWeight: 500 }}>{lb}</span>{ct&&<span style={{ fontSize: 14, color: C.muted, fontWeight: 600 }}>{ct}</span>}<span style={{ color: C.muted }}>›</span>
          </button>
        ))}
      </div>
    </div>
  );
}

// ── CHAT LIST ──
function ChatList({ nav }) {
  const chats = [
    { id: 1, name: "Алексей Маратов", av: "АМ", msg: "Да, best of 3 давай 💪", time: "14:39", unread: 2, type: "p", on: true },
    { id: 2, name: "NTC Astana", av: "NT", msg: "📢 Обновление расписания...", time: "12:05", unread: 15, type: "c" },
    { id: 3, name: "Парная игра 10 мар", av: "🎾", msg: "Все подтвердили!", time: "Вчера", unread: 0, type: "e" },
    { id: 4, name: "Сауле Абдрахманова", av: "СА", msg: "Спасибо за игру!", time: "Вчера", unread: 0, type: "p" },
  ];
  return (
    <div style={{ flex: 1, overflow: "auto" }}>
      {chats.map(c=>(
        <button key={c.id} onClick={()=>nav("chatDetail")} style={{ display: "flex", alignItems: "center", gap: 12, padding: "14px 16px", width: "100%", border: "none", borderBottom: `1px solid ${C.borderL}`, background: c.unread>0?C.primaryLight+"40":C.card, cursor: "pointer", textAlign: "left" }}>
          {c.type==="e"?<div style={{ width: 48, height: 48, borderRadius: 24, background: C.primaryLight, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 22 }}>{c.av}</div>:<Av text={c.av} size={48} online={c.on} color={c.type==="c"?C.primary:undefined}/>}
          <div style={{ flex: 1, minWidth: 0 }}>
            <div style={{ display: "flex", justifyContent: "space-between" }}>
              <div style={{ display: "flex", alignItems: "center", gap: 4 }}><span style={{ fontWeight: c.unread>0?700:500, fontSize: 15 }}>{c.name}</span>{c.type==="c"&&<span style={{ color: C.primary, fontSize: 12 }}>✓</span>}{c.type==="e"&&<span style={{ fontSize: 10, color: C.blue, background: C.blue+"15", padding: "1px 6px", borderRadius: 4, fontWeight: 600 }}>Ивент</span>}</div>
              <span style={{ fontSize: 12, color: C.muted }}>{c.time}</span>
            </div>
            <div style={{ display: "flex", justifyContent: "space-between", marginTop: 2 }}>
              <span style={{ fontSize: 14, color: c.unread>0?C.text:C.muted, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap", maxWidth: 220 }}>{c.msg}</span>
              {c.unread>0&&<span style={{ minWidth: 20, height: 20, borderRadius: 10, background: C.primary, color: "white", fontSize: 11, fontWeight: 700, display: "flex", alignItems: "center", justifyContent: "center", padding: "0 6px" }}>{c.unread}</span>}
            </div>
          </div>
        </button>
      ))}
    </div>
  );
}

// ── CHAT DETAIL ──
function ChatDetail() {
  const [msg, setMsg] = useState("");
  return (
    <div style={{ flex: 1, display: "flex", flexDirection: "column" }}>
      <div style={{ flex: 1, overflow: "auto", padding: 16, display: "flex", flexDirection: "column", gap: 8 }}>
        {MSGS.map(m=>(
          <div key={m.id} style={{ display: "flex", justifyContent: m.me?"flex-end":"flex-start" }}>
            <div style={{ maxWidth: "75%", padding: "10px 14px", borderRadius: m.me?"18px 18px 4px 18px":"18px 18px 18px 4px", background: m.me?C.primary:C.card, color: m.me?"white":C.text, border: m.me?"none":`1px solid ${C.borderL}` }}>
              <p style={{ fontSize: 15, margin: 0, lineHeight: 1.4 }}>{m.text}</p>
              <p style={{ fontSize: 11, margin: "4px 0 0", opacity: 0.6, textAlign: "right" }}>{m.t}{m.me&&" ✓✓"}</p>
            </div>
          </div>
        ))}
      </div>
      <div style={{ padding: "10px 16px 28px", background: C.card, borderTop: `1px solid ${C.borderL}`, display: "flex", gap: 10, alignItems: "flex-end" }}>
        <input value={msg} onChange={e=>setMsg(e.target.value)} placeholder="Сообщение..." style={{ flex: 1, padding: "12px 16px", borderRadius: 24, border: `1px solid ${C.border}`, fontSize: 15, outline: "none" }}/>
        <button style={{ width: 44, height: 44, borderRadius: 22, background: msg?C.primary:C.borderL, border: "none", color: msg?"white":C.muted, fontSize: 18, cursor: "pointer", display: "flex", alignItems: "center", justifyContent: "center" }}>↑</button>
      </div>
    </div>
  );
}

// ── NOTIFICATIONS ──
function NotifsScreen() {
  const notifs = [
    { icon: "🎾", title: "Подтвердите результат", body: "Алексей внёс: 6-4, 6-3", time: "2ч", unread: true },
    { icon: "📊", title: "Рейтинг вырос!", body: "Позиция #45 (+18.5)", time: "3ч", unread: true },
    { icon: "💬", title: "Новое сообщение", body: "Алексей: Да, best of 3 💪", time: "5ч", unread: true },
    { icon: "🏆", title: "Новое достижение!", body: 'Вы получили бейдж "Десятка"', time: "1д", unread: false },
    { icon: "🎾", title: "Напоминание", body: "Парная игра завтра в 10:00", time: "1д", unread: false },
    { icon: "👥", title: "Заявка одобрена", body: "Вы вступили в NTC Astana", time: "2д", unread: false },
  ];
  return (
    <div style={{ flex: 1, overflow: "auto" }}>
      <p style={{ padding: "12px 16px 4px", fontSize: 13, fontWeight: 700, color: C.muted }}>Сегодня</p>
      {notifs.slice(0,3).map((n,i)=>(
        <div key={i} style={{ display: "flex", gap: 12, padding: "12px 16px", background: n.unread?C.primaryLight+"40":C.card, borderBottom: `1px solid ${C.borderL}` }}>
          <div style={{ width: 40, height: 40, borderRadius: 20, background: C.primaryLight, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 18, flexShrink: 0 }}>{n.icon}</div>
          <div style={{ flex: 1 }}><div style={{ display: "flex", justifyContent: "space-between" }}><span style={{ fontWeight: 600, fontSize: 14 }}>{n.title}</span><span style={{ fontSize: 12, color: C.muted }}>{n.time}</span></div><p style={{ fontSize: 13, color: C.sub, margin: "2px 0 0" }}>{n.body}</p></div>
          {n.unread&&<div style={{ width: 8, height: 8, borderRadius: 4, background: C.primary, flexShrink: 0, marginTop: 6 }}/>}
        </div>
      ))}
      <p style={{ padding: "12px 16px 4px", fontSize: 13, fontWeight: 700, color: C.muted }}>Ранее</p>
      {notifs.slice(3).map((n,i)=>(
        <div key={i} style={{ display: "flex", gap: 12, padding: "12px 16px", background: C.card, borderBottom: `1px solid ${C.borderL}` }}>
          <div style={{ width: 40, height: 40, borderRadius: 20, background: C.borderL, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 18, flexShrink: 0 }}>{n.icon}</div>
          <div style={{ flex: 1 }}><div style={{ display: "flex", justifyContent: "space-between" }}><span style={{ fontWeight: 600, fontSize: 14 }}>{n.title}</span><span style={{ fontSize: 12, color: C.muted }}>{n.time}</span></div><p style={{ fontSize: 13, color: C.sub, margin: "2px 0 0" }}>{n.body}</p></div>
        </div>
      ))}
    </div>
  );
}

// ── PLAYER PROFILE ──
function PlayerProfile() {
  return (
    <div style={{ flex: 1, overflow: "auto" }}>
      <div style={{ padding: "24px 16px", background: C.card, textAlign: "center" }}>
        <Av text="МК" size={72}/>
        <h2 style={{ fontSize: 22, fontWeight: 800, margin: "12px 0 4px" }}>Марат Касымов</h2>
        <p style={{ fontSize: 14, color: C.sub, margin: 0 }}>🎾 Продвинутый 4.5 · Есильский</p>
        <div style={{ display: "flex", justifyContent: "center", gap: 10, marginTop: 16 }}>
          <button style={{ padding: "10px 20px", borderRadius: 12, background: C.primary, color: "white", border: "none", fontWeight: 600, cursor: "pointer" }}>✉️ Написать</button>
          <button style={{ padding: "10px 20px", borderRadius: 12, background: C.primaryLight, color: C.primary, border: `1px solid ${C.primary}`, fontWeight: 600, cursor: "pointer" }}>🎾 Пригласить</button>
          <button style={{ padding: "10px 20px", borderRadius: 12, background: C.borderL, border: "none", fontWeight: 600, cursor: "pointer" }}>👥</button>
        </div>
      </div>
      <div style={{ display: "flex", justifyContent: "space-around", padding: 16, background: C.card, borderTop: `1px solid ${C.borderL}` }}>
        {[["120","Игр",C.primary],["72%","Побед",C.green],["1,650","Рейтинг",C.accent],["#1","NTC",C.text]].map(([v,l,c])=><div key={l} style={{ textAlign: "center" }}><p style={{ fontSize: 20, fontWeight: 800, color: c, margin: 0 }}>{v}</p><p style={{ fontSize: 11, color: C.muted, margin: 0 }}>{l}</p></div>)}
      </div>
      <div style={{ padding: 16 }}>
        <span style={{ fontWeight: 700, fontSize: 15 }}>Достижения</span>
        <div style={{ display: "flex", gap: 8, marginTop: 10, overflowX: "auto" }}>
          {BADGES.filter(b=>b.earned).map(b=><div key={b.id} style={{ textAlign: "center", minWidth: 60 }}><div style={{ width: 44, height: 44, borderRadius: 12, background: C.primaryLight, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 22, margin: "0 auto" }}>{b.icon}</div><p style={{ fontSize: 10, margin: "4px 0 0" }}>{b.name}</p></div>)}
        </div>
      </div>
      <div style={{ padding: "0 16px" }}>
        <span style={{ fontWeight: 700, fontSize: 15 }}>Сообщества</span>
        <div style={{ display: "flex", gap: 8, marginTop: 10 }}>
          {["NTC Astana","Astana League"].map(n=><div key={n} style={{ display: "flex", alignItems: "center", gap: 6, padding: "8px 14px", borderRadius: 12, background: C.borderL }}><Av text={n.slice(0,2)} size={24}/><span style={{ fontSize: 13, fontWeight: 500 }}>{n}</span></div>)}
        </div>
      </div>
    </div>
  );
}

// ── COMMUNITIES LIST ──
function CommunitiesList({ nav }) {
  return (
    <div style={{ flex: 1, overflow: "auto", padding: 16 }}>
      <div style={{ display: "flex", alignItems: "center", gap: 8, padding: "10px 14px", borderRadius: 12, background: C.card, border: `1px solid ${C.border}`, marginBottom: 12 }}>
        <span style={{ color: C.muted }}>🔍</span><input placeholder="Поиск сообществ..." style={{ flex: 1, border: "none", outline: "none", fontSize: 15, background: "transparent" }}/>
      </div>
      <div style={{ display: "flex", gap: 8, marginBottom: 12, overflowX: "auto" }}>{["Все","Клубы","Лиги","Группы"].map((f,i)=><Chip key={f} label={f} active={i===0}/>)}</div>
      {[
        { name: "NTC Astana", type: "Клуб", mem: 245, v: true },
        { name: "Astana Tennis League", type: "Лига", mem: 180, v: true },
        { name: "Weekend Tennis", type: "Группа", mem: 45, v: false },
        { name: "Pro Tennis KZ", type: "Организатор", mem: 320, v: true },
      ].map((c,i)=>(
        <Cd key={i} onClick={()=>nav("communityDetail")} style={{ marginBottom: 10, display: "flex", alignItems: "center", gap: 12 }}>
          <div style={{ width: 48, height: 48, borderRadius: 14, background: C.primaryLight, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 18, fontWeight: 700, color: C.primary }}>{c.name.slice(0,2)}</div>
          <div style={{ flex: 1 }}>
            <div style={{ display: "flex", alignItems: "center", gap: 4 }}><span style={{ fontWeight: 700, fontSize: 15 }}>{c.name}</span>{c.v&&<span style={{ color: C.primary, fontSize: 14 }}>✓</span>}</div>
            <p style={{ fontSize: 12, color: C.muted, margin: "2px 0" }}>{c.type} · {c.mem} участников</p>
          </div>
        </Cd>
      ))}
    </div>
  );
}

// ── MAIN APP ──
export default function App() {
  const [scr, setScr] = useState("auth_phone");
  const [tab, setTab] = useState("home");
  const [hist, setHist] = useState([]);

  const nav = s => { setHist(p=>[...p,scr]); setScr(s); };
  const back = () => { const p = hist[hist.length-1]||"home"; setHist(h=>h.slice(0,-1)); setScr(p); };
  const isMain = ["home","players","events","communities","profile"].includes(scr);
  useEffect(() => { if(isMain) setTab(scr); }, [scr]);
  const tabChange = t => { setTab(t); setScr(t); setHist([]); };

  return (
    <div style={{ minHeight: "100vh", display: "flex", alignItems: "center", justifyContent: "center", background: "linear-gradient(135deg,#e8f5ee 0%,#f0f9ff 50%,#fef3c7 100%)", padding: "20px 0" }}>
      <Frame>
        <SBar/>
        {scr==="auth_phone"&&<AuthPhone onNext={()=>nav("auth_otp")}/>}
        {scr==="auth_otp"&&<AuthOTP onNext={()=>nav("auth_profile")} onBack={back}/>}
        {scr==="auth_profile"&&<AuthProfile onNext={()=>{setScr("home");setTab("home");setHist([]);}}/>}

        {isMain&&<>
          <Header title={{home:"Tennis Astana",players:"Игроки",events:"Ивенты",communities:"Сообщества",profile:"Профиль"}[scr]} right={<HeaderBadges navigate={nav}/>}/>
          {scr==="home"&&<HomeScreen nav={nav}/>}
          {scr==="players"&&<PlayersScreen nav={nav}/>}
          {scr==="events"&&<EventsScreen nav={nav}/>}
          {scr==="communities"&&<CommunitiesList nav={nav}/>}
          {scr==="profile"&&<ProfileScreen nav={nav}/>}
          <TabBar active={tab} onChange={tabChange}/>
        </>}

        {scr==="eventDetail"&&<><Header title="Детали ивента" onBack={back}/><EventDetail onBack={back}/></>}
        {scr==="eventCreate"&&<><Header title="Новый ивент" onBack={back}/><EventCreate onBack={back}/></>}
        {scr==="communityDetail"&&<><Header title="" onBack={back} right={<button style={{ background: "none", border: "none", fontSize: 18, cursor: "pointer" }}>⋯</button>}/><CommunityScreen onBack={back}/></>}
        {scr==="chatList"&&<><Header title="Чаты" onBack={back}/><ChatList nav={nav}/></>}
        {scr==="chatDetail"&&<><Header title="Алексей Маратов" onBack={back} right={<span style={{ fontSize: 12, color: C.green }}>● онлайн</span>}/><ChatDetail/></>}
        {scr==="playerProfile"&&<><Header title="Профиль" onBack={back}/><PlayerProfile/></>}
        {scr==="notifs"&&<><Header title="Уведомления" onBack={back}/><NotifsScreen/></>}
      </Frame>
    </div>
  );
}
