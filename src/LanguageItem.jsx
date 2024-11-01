import React from 'react';
import styles from './LanguageItem.module.css';

function LanguageItem({ rank, language }) {
  const logoPath = `./logos/${language.name.toLowerCase()}.svg`;

  return (
    <div className={styles.item}>
      <span className={styles.rank}>#{rank}</span>
      <img src={logoPath} alt={`${language.name} logo`} className={styles.logo} />
      <a href={language.url} target="_blank" rel="noopener noreferrer" className={styles.name}>
        {language.name}
      </a>
      <span className={styles.stars}>{language.stars.toLocaleString()} ⭐</span>
    </div>
  );
}

export default LanguageItem;
