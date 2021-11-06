import type { NextPage } from "next";
import Head from "next/head";
import Image from "next/image";
import styles from "../styles/Home.module.css";

import Distance from "../components/Distance";
import Container from "@mui/material/Container";

const Home: NextPage = () => {
  return (
    <Container maxWidth="lg">
      <h1>Nazdik</h1>
      <h4>How close are you?</h4>
      <Distance />
    </Container>
  );
};

export default Home;
