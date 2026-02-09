<template>
  <div class="article-body">
    <!-- Table of Contents -->
    <nav class="toc">
      <h3 class="toc-title">Table of Contents</h3>
      <ol class="toc-list">
        <li>
          <a href="#introduction" class="toc-link">Introduction</a>
          <ul>
            <li><a href="#why-fortran" class="toc-link">Why Fortran?</a></li>
          </ul>
        </li>
        <li><a href="#tech-stack" class="toc-link">Tech Stack</a></li>
        <li>
          <a href="#initial-simulation" class="toc-link"
            >The Charged Particle in a Magnetic Field</a
          >
        </li>
        <li>
          <a href="#implementation" class="toc-link">Implementation</a>
          <ul>
            <li><a href="#the-simulation" class="toc-link">The Simulation</a></li>
            <li><a href="#the-interface" class="toc-link">The Interface</a></li>
            <li><a href="#the-execution" class="toc-link">The Execution</a></li>
          </ul>
        </li>
        <li><a href="#results" class="toc-link">Results</a></li>
        <li><a href="#outcomes-and-learnings" class="toc-link">Outcomes and Learnings</a></li>
      </ol>
    </nav>

    <!-- Introduction -->
    <h2 id="introduction" class="section-heading">Introduction</h2>
    <p>
      I've been a software engineer for quite some time now. Like many developers, the AI wave has
      become something utterly inescapable. Regardless of where I look, what I work on, or what I
      read, I can't get more than 10 minutes in before the topic of AI comes up.
    </p>

    <p>
      Now don't get me wrong, I'm not anti AI. Staying relevant in tech requires a degree of
      adaptation. Agentic development will become a major part of how we work in software over the
      next decade, of that there can be little doubt. So I do what I suggest to anyone who will
      listen: make an effort to get good at using AI. That being said, I find myself needing some
      downtime from it on a daily basis.
    </p>

    <p>
      There are also real dangers to using AI for everything. Problem solving is a skill that needs
      to be actively developed and practiced. The temptation to ask an LLM for assistance on every
      little problem can become overwhelming, and we are in real danger of creating a generation of
      people that are incapable of problem solving. Your brain is a muscle like any other, it needs
      to be worked or it will deteriorate.
    </p>

    <p>
      Enter the simulation catalogue. This application/codebase is my down time from AI. At its
      heart, it's a Fortran physics engine that is indexed and made available via this app. All of
      it is written without AI. No code completions, no Claude Code. Nothing but pen and paper, a
      basic IDE and a few books. Distractions are turned off and U just focus on the mathematics and
      the code.
    </p>
    <p>
      I'm a firm believer that solving difficult problems is an essential part of what it means to
      be human. Much like a sheep dog that requires a job to be happy, difficult problems provide a
      positive channel for our internal machinations. I sleep much better if I sit down after my
      kids are asleep and work on a few physics problems. Maybe that's just me.
    </p>

    <h3 id="why-fortran" class="subsection-heading">Why Fortran?</h3>
    <p>
      Many people have asked: why Fortran? Why not just use NumPy? It's a valid question, and I have
      three good reasons to pick Fortran over Python.
    </p>

    <ol>
      <li>
        Performance - NumPy is impressively fast with vectorized operations. Still not as fast as
        Fortran, but not far off. However, many problems in physics can't be vectorized (Monte Carlo
        simulations are probably the most famous example of this). Sometimes, a good old fashioned
        loop is the only option, and Fortran is several orders of magnitude faster in these cases.
      </li>
      <li class="q-mt-md">
        Marketability - Fortran is still very much a marketable skill, one that's becoming
        increasingly rare. You'd be amazed at the amount of critical infrastructure code that is
        written in Fortran. It may be a niche market, but with the advent of AI, developing some
        niche skills is likely not a bad thing.
      </li>
      <li class="q-mt-md">
        Variety - I've been working on Python web apps since I started working in web development.
        Working with Fortran provides a nice change of scenery.
      </li>
    </ol>

    <p>
      It's worth mentioning that I've worked with Fortran extensively in the past. I have a degree
      in Theoretical Physics, and much of the computational work that I did was in Fortran. It was
      also the first language I learned to program in, so it offers a certain sense of nostalgia.
    </p>

    <!-- Tech Stack -->
    <h2 id="tech-stack" class="section-heading">Tech Stack</h2>
    <p>
      Now that the boring stuff is out of the way, let's focus on the engineering. I stand on the
      shoulders of giants; while the aim of this project is to work on being independent on external
      tools, there are a few I can't do without.
    </p>

    <div class="tech-stack-grid">
      <!-- Fortran -->
      <div class="tech-card">
        <div class="tech-icon">
          <img
            src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/fortran/fortran-original.svg"
            alt="Fortran"
            width="48"
            height="48"
          />
        </div>
        <h4 class="tech-title">Fortran</h4>
        <p class="tech-description">
          The backbone of the physics engine. When it comes to numerical computing, Fortran is still
          king. It's incredibly quick, and much of what you need to work on physics problems is
          native to Fortran.
        </p>
        <div class="tech-tags">
          <span class="tech-tag">simulations</span>
        </div>
      </div>

      <!-- Golang -->
      <div class="tech-card">
        <div class="tech-icon">
          <img
            src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original-wordmark.svg"
            alt="Golang"
            width="48"
            height="48"
          />
        </div>
        <h4 class="tech-title">Golang</h4>
        <p class="tech-description">
          My go-to programming language. Unless I have a good reason to use something else, I use
          Golang. It's performant, has a ton of industry support, and provides an unmatched balance
          between development speed and code stability.
        </p>
        <div class="tech-tags">
          <span class="tech-tag">api</span>
          <span class="tech-tag">worker</span>
        </div>
      </div>

      <!-- Terraform -->
      <div class="tech-card">
        <div class="tech-icon">
          <img
            src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/terraform/terraform-original.svg"
            alt="Terraform"
            width="48"
            height="48"
          />
        </div>
        <h4 class="tech-title">Terraform</h4>
        <p class="tech-description">
          Any infrastructure I need is always provisioned with Terraform. Its the only IaC solution
          I've ever needed, and the sheer number of providers it offers makes it a no brainer.
        </p>
        <div class="tech-tags">
          <span class="tech-tag">infrastructure</span>
          <span class="tech-tag">ci/cd</span>
        </div>
      </div>

      <!-- Vue -->
      <div class="tech-card">
        <div class="tech-icon">
          <img
            src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/vuejs/vuejs-original.svg"
            alt="Vue"
            width="48"
            height="48"
          />
        </div>
        <h4 class="tech-title">Vue</h4>
        <p class="tech-description">
          I'm not much of a FE engineer, but Vue has always been my go to. Many people will say
          "Just use React", and maybe they're right.
        </p>
        <div class="tech-tags">
          <span class="tech-tag">frontend</span>
        </div>
      </div>

      <!-- Kubernetes -->
      <div class="tech-card">
        <div class="tech-icon">
          <img
            src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/kubernetes/kubernetes-plain.svg"
            alt="Kubernetes"
            width="48"
            height="48"
          />
        </div>
        <h4 class="tech-title">Kubernetes</h4>
        <p class="tech-description">
          All services are deployed on a Kubernetes cluster that I maintain on a Debian virtual
          machine somewhere in the cloud.
        </p>
        <div class="tech-tags">
          <span class="tech-tag">orchestration</span>
        </div>
      </div>

      <!-- PostgreSQL -->
      <div class="tech-card">
        <div class="tech-icon">
          <img
            src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-original.svg"
            alt="PostgreSQL"
            width="48"
            height="48"
          />
        </div>
        <h4 class="tech-title">PostgreSQL</h4>
        <p class="tech-description">
          PostgreSQL handles all persistent data storage. My Kubernetes cluster already features a
          fully operational PostgreSQL cluster so this was a no brainer.
        </p>
        <div class="tech-tags">
          <span class="tech-tag">database</span>
        </div>
      </div>
    </div>

    <!-- Initial Simulation -->
    <h2 id="initial-simulation" class="section-heading">
      The Charged Particle in a Magnetic Field
    </h2>
    <p>
      There are a lot of problems I could have picked. In the end, I decided to go for something
      simple. A charged particle moving in a constant electro-magnetic field, with the addition of
      gravity.
    </p>

    <div class="equation" v-html="forceEquation"></div>

    <p>
      This can easily be solved analytically. However, with the amount of engineering work that I
      had to do, I wanted to start with a problem that had a well defined solution, and that I was
      thoroughly familiar with. The equations of motion become:
    </p>

    <div class="equation" v-html="equationOfMotion"></div>

    <p>
      where <span v-html="renderLatexInline(`\\omega_c`)"></span> is the cyclotron frequency and
      <span v-html="renderLatexInline(`\\mathbf{v}_D `)"></span> is the drift velocity. A charged
      particle moving in this system will create a helical trajectory. The addition of gravity
      introduces a horizontal drift that can be counter-balanced with an appropriate electric field.
    </p>

    <!-- Implementation -->
    <h2 id="implementation" class="section-heading">Implementation</h2>
    <p>Since this was the first simulation I wrote, there are three distinct problems to solve:</p>

    <ol>
      <li>Writing the simulation - handled in Fortran</li>
      <li>Running the simulation - handled in Golang</li>
      <li>Displaying the results - combination of Golang and Vue</li>
    </ol>

    <h3 id="the-simulation" class="subsection-heading">The Simulation</h3>
    <p>
      The first part was the easiest, and took the least amount of time. It took me a few hours to
      get back up to speed with Fortran, but nothing major to report. The initial code featured a
      velocity verlet algorithm, but this produced unstable trajectories and was quickly replaced
      with a Boris Pusher integrator (more on that in a future article).
    </p>

    <h3 id="the-interface" class="subsection-heading">The Interface</h3>
    <p>
      In order to actually run the simulation, the Fortran code needed to be able to interface with
      the Golang components. This turned out to be the most difficult part. While technically not a
      difficult problem to solve for one simulation, defining an interface that works for ANY
      simulation in a clean, compact and repeatable manner is a completely different problem.
    </p>
    <p>
      In the end, I settled on a fairly simple solution: TOML configuration files. All simulations
      were written in such a way that the first and only argument was the path to a TOML
      configuration file. This file would contain all the parameters and configuration settings
      required to run the simulation and produce any output artifacts. An example config looks like
      the following.
    </p>

    <div class="code-block" v-html="highlightCode(exampleTomlCode, 'toml')"></div>
    <p>
      On startup, each Fortran simulation would load the provided TOML file, and initialize all
      defined variables using the data in the file. The simulation would then run, and all outputs
      were saved in CSV format to an output folder specified in the config file. In (rough) Fortran:
    </p>

    <div class="code-block" v-html="highlightCode(exampleFortranCodeRaw, 'fortran')"></div>

    <h3 id="the-execution" class="subsection-heading">The Execution</h3>
    <p>
      With the interface in place, running the simulation was straightforward. I compiled all the
      Fortran code, and stored in blob storage, along with relevant metadata that identified the
      simulation and any required parameters. When the API received a request from a client
      containing a unique simulation ID and input parameters, it sent a message down a job queue,
      and a worker listening on the other end runs the simulation with the following steps:
    </p>

    <ol>
      <li>Download relevant simulation executable from blob storage</li>
      <li>
        Create a temporary directory containing the compiled simulation executable and output
        directory
      </li>
      <li>
        Write the user provided simulation parameters to a TOML file in the temporary directory
      </li>
      <li>Run the simulation executable</li>
      <li>Zip the contents of the output directory (populated with CSV files by the simulation)</li>
      <li>Store output zip in blob storage</li>
    </ol>

    <p>
      This flow provides a simple framework that can be used to run any simulation, and it's all
      based on two simple rules: TOML in, CSV out.
    </p>

    <p>
      The final piece of the puzzle was displaying the results. This was done with a combination of
      Vue and Quasar. I used
      <a href="https://plotly.com/javascript/" target="_blank" class="inline-link">Plotly.js</a> to
      actually display the results. Much like any self-respecting backend/infrastructure engineer,
      this is the part that I enjoyed the least.
    </p>

    <h2 id="results" class="section-heading">Results</h2>

    <p>
      With everything in place, it was time to run some simulations. The following plots show the
      results of a charged particle moving through a uniform magnetic field.
      <span v-html="renderLatexInline(`x_1`)"></span>,
      <span v-html="renderLatexInline(`x_2`)"></span> and
      <span v-html="renderLatexInline(`x_3`)"></span> are used for the basis, and in this case
      correspond with cartesian coordinates <span v-html="renderLatexInline(`x`)"></span>,
      <span v-html="renderLatexInline(`y`)"></span>, <span v-html="renderLatexInline(`z`)"></span>.
    </p>

    <div class="results-gallery">
      <figure class="result-figure">
        <img
          src="/blog/images/001/helical_motion_3d.png"
          alt="3D Helical Motion Trajectory"
          class="result-image"
        />
        <figcaption class="result-caption">
          Figure 1 <br />
          3D trajectory of a charged particle in a magnetic field, showing the characteristic
          helical motion.
        </figcaption>
      </figure>

      <figure class="result-figure">
        <img
          src="/blog/images/001/position_evolution.png"
          alt="Position Evolution Over Time"
          class="result-image"
        />
        <figcaption class="result-caption">
          Figure 2 <br />
          Time evolution of the particle's position components.
        </figcaption>
      </figure>
    </div>

    <p>
      Figure 1 shows the helical motion created in the presence of a magnetic field. Note that the
      helix is compressed in the <span v-html="renderLatexInline(`x_2`)"></span> plane due to the
      addition of gravity.
    </p>

    <p>
      Figure 2 shows the time evolution of the individual position components. The magnetic field
      was set to be aligned with the <span v-html="renderLatexInline(`x_1`)"></span> axis, resulting
      in a constant velocity in the <span v-html="renderLatexInline(`x_1`)"></span> plane. Both the
      <span v-html="renderLatexInline(`x_2`)"></span> and
      <span v-html="renderLatexInline(`x_3`)"></span> planes experience acceleration due to the
      magnetic field, producing the classic helical trajectory.
    </p>

    <p>Figure 2 also clearly shows the horizontal drift created by the introduction of gravity.</p>

    <!-- Outcomes and Learnings -->
    <h2 id="outcomes-and-learnings" class="section-heading">Outcomes and Learnings</h2>
    <p>
      I set out to find a useful way to spend some down time and get away from the AI craze. And it
      has to be said, Fortran and Physics delivered.
    </p>
    <p>
      The web development portion is what I do everyday, and these are arguably the least
      interesting components of the project. However, I'm happy with the overall framework and how
      it came together. The most important thing was to develop something that was generic so that I
      could in future focus on what really interested me: the physics and the Fortran. And I think
      to that degree, this initial implementation was a success. The current web app can run any
      Fortran simulation, provided it follows the TOML in, CSV out convention.
    </p>
    <p>In regards to the Fortran and the physics, I learned two things:</p>
    <ol>
      <li>Keep the Fortran minimal</li>
      <li>Don't waste your time with Verlet Velocity in a magnetic field</li>
    </ol>

    <p>
      Speaking to the first point, Fortran is not a general purpose language. It's hyper-optimized
      for numerical computing, and it's incredibly good at that. However, everything else can be a
      challenge. Keeping the Fortran code scoped to (1) read config, (2) run simulation, (3) output
      CSV, solved many headaches, and produced a clean, and maintainable codebase. The motto from
      there on very much became: let Fortran do what it's built for, do everything else in Golang.
    </p>

    <p>
      As for the second point, I can't really speak ill of Verlet Velocity. It's incredibly easy to
      implement and works well for a whole host of systems. However, when magnetism is involved,
      just stick with Boris. It wasn't much harder to implement and produced far, far better
      results. This won't come as a surprise to anyone experienced in the field, but for a lowly web
      developer, it was quite the revelation.
    </p>

    <p>
      That's it for today. For anyone who's interested, the entire codebase is available on my
      public GitHub:
      <a
        href="https://github.com/psauerborn/simulation-catalogue"
        target="_blank"
        class="inline-link"
        >simulation-catalogue</a
      >
      (API, worker, and web app) and
      <a href="https://github.com/psauerborn/simulations" target="_blank" class="inline-link"
        >simulations</a
      >
      (Fortran physics engine).
    </p>
  </div>
</template>

<script setup>
import { highlightCode, renderLatexInline, renderLatexDisplay } from './utils.js'

// force equation for lorentz force + gravity
const forceEquation = renderLatexDisplay(
  `\\mathbf{F} = q(\\mathbf{E} + \\mathbf{v} \\times \\mathbf{B}) + m\\mathbf{g}`,
)

const equationOfMotion = renderLatexDisplay(
  `\\begin{aligned}
  m\\frac{d\\mathbf{v}}{dt} &= q(\\mathbf{E} + \\mathbf{v} \\times \\mathbf{B}) + m\\mathbf{g} \\\\[1em]
  \\omega_c &= \\frac{qB}{m} \\\\[1em]
  \\mathbf{v}_D &= \\frac{\\mathbf{E} \\times \\mathbf{B}}{B^2}
  \\end{aligned}`,
)

const exampleTomlCode = `[parameters]
magnetic_field = [0.0, 0.0, 1.0]
electric_field = [0.0, 0.0, 0.0]

charge = 1.0
mass = 1.0

initial_position = [0.0, 0.0, 0.0]
initial_velocity = [1.0, 0.0, 0.0]

delta_t = 0.01
num_steps = 5000

[configuration]
output_dir = output`

const exampleFortranCodeRaw = `program simulation
  implicit none

  real, dimension(3) :: parameter_1
  real :: parameter_2

  character(len=:), allocatable :: output_dir
  character(len=256) :: output_path

  real, dimension(10000, 3) :: output
  integer :: i

  call load_config()

  do i=1, 10000
    ! execute simulation here
    output(i, :) = [1.0, 2.0, 3.0]
  end do

  output_path = trim(output_dir)//'/'//'output.csv'

  call write_csv(output_path, output)


contains
  subroutine load_config()
      character(len=256) :: config_path = "etc/default_config.toml"
      integer :: n_args, stat
      type(toml_table), allocatable :: table
      type(toml_error), allocatable :: error

      n_args = command_argument_count()
      if (n_args > 0) then
         call get_command_argument(1, config_path)
      end if

      call toml_load(table, config_path, error=error)
      if (allocated(error)) then
         print '(a)', error%message
         stop 1
      end if

      call get_value(table, "parameter_1", parameter_1, stat=stat)
      call get_value(table, "parameter_2", parameter_2, stat=stat)
      call get_value(table, "output_dir", output_dir, stat=stat)

  end subroutine load_config
end program`
</script>

<style lang="scss" scoped>
.toc {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  padding: 20px 24px;
  margin-bottom: 32px;

  .toc-title {
    font-size: 0.85rem;
    font-weight: 600;
    color: #636366;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    margin: 0 0 12px 0;
  }

  .toc-list {
    margin: 0;
    padding-left: 20px;
    color: #d1d1d6;

    li {
      margin-bottom: 6px;
      line-height: 1.5;
    }

    ul {
      padding-left: 20px;
      margin: 6px 0 0 0;
    }
  }

  .toc-link {
    color: #30d158;
    text-decoration: none;
    font-size: 0.95rem;

    &:hover {
      text-decoration: underline;
    }
  }
}

.section-heading {
  font-size: 1.5rem;
  font-weight: 700;
  color: #f5f5f7;
  margin: 40px 0 16px 0;
  letter-spacing: -0.01em;

  &:first-of-type {
    margin-top: 0;
  }
}

.subsection-heading {
  font-size: 1.2rem;
  font-weight: 600;
  color: #e5e5ea;
  margin: 28px 0 12px 0;
}

.equation {
  margin: 24px 0;
  padding: 20px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  text-align: center;
  overflow-x: auto;
}

.tech-stack-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  margin: 24px 0 32px 0;
}

.tech-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  padding: 24px;
  transition: all 0.3s ease;

  &:hover {
    background: rgba(255, 255, 255, 0.06);
    border-color: rgba(255, 255, 255, 0.15);
    transform: translateY(-2px);
  }
}

.tech-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  margin-bottom: 16px;

  svg {
    width: 32px;
    height: 32px;
  }
}

.tech-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #f5f5f7;
  margin: 0 0 8px 0;
}

.tech-description {
  font-size: 0.9rem;
  line-height: 1.6;
  color: #a1a1a6;
  margin: 0;
}

.tech-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.tech-tag {
  display: inline-block;
  font-size: 0.75rem;
  font-weight: 500;
  color: #30d158;
  background: rgba(48, 209, 88, 0.1);
  border: 1px solid rgba(48, 209, 88, 0.25);
  border-radius: 6px;
  padding: 4px 10px;
  letter-spacing: 0.02em;
}

.code-block {
  margin: 24px 0;

  :deep(pre) {
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 8px;
    padding: 20px;
    overflow-x: auto;
    margin: 0;

    code {
      font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
      font-size: 0.875rem;
      line-height: 1.6;
      color: #d1d1d6;
    }
  }
}

.results-gallery {
  display: flex;
  flex-direction: column;
  gap: 24px;
  margin: 24px 0 32px 0;
}

.result-figure {
  margin: 0;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;

  &:hover {
    background: rgba(255, 255, 255, 0.06);
    border-color: rgba(255, 255, 255, 0.15);
    transform: translateY(-2px);
  }
}

.result-image {
  width: 100%;
  height: auto;
  display: block;
}

.result-caption {
  padding: 16px;
  font-size: 0.9rem;
  line-height: 1.5;
  color: #a1a1a6;
  text-align: center;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}
</style>
