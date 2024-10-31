import React from 'react';
import styles from './LanguageItem.module.css';

function LanguageItem({ rank, language }) {
  return (
    <div className={styles.item}>
      <span className={styles.rank}>#{rank}</span>
      <a href={language.url} target="_blank" rel="noopener noreferrer" className={styles.name}>
        {language.name}
      </a>
      <span className={styles.stars}>{language.stars.toLocaleString()} ⭐</span>
    </div>
  );
}

export default LanguageItem;
