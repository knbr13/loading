import React from 'react';
import LanguageItem from '../LanguageItem';
import styles from './LanguageList.module.css';

function LanguageList({ languages }) {
  return (
    <div className={styles.list}>
      {languages
        .sort((a, b) => b.stars - a.stars)  // Sort by stars
        .map((lang, index) => (
          <LanguageItem key={lang.name} rank={index + 1} language={lang} />
        ))}
    </div>
  );
}

export default LanguageList;
